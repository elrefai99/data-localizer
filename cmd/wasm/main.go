//go:build js && wasm

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"syscall/js"

	localizer "github.com/elrefai99/data-localizer"
)

type request struct {
	Data             json.RawMessage `json:"data"`
	LanguageHeader   string          `json:"languageHeader"`
	FallbackLanguage string          `json:"fallbackLanguage"`
	Options          *requestOptions `json:"options"`
}

type requestOptions struct {
	SupportedLanguages []string `json:"supportedLanguages"`
	FallbackLanguage   string   `json:"fallbackLanguage"`
	FallbackLanguages  []string `json:"fallbackLanguages"`
	MissingTranslation string   `json:"missingTranslation"`
}

type response struct {
	OK    bool           `json:"ok"`
	Value any            `json:"value"`
	Error *responseError `json:"error,omitempty"`
}

type responseError struct {
	Name               string   `json:"name"`
	Message            string   `json:"message"`
	Path               string   `json:"path,omitempty"`
	OriginalPath       string   `json:"originalPath,omitempty"`
	AttemptedLanguages []string `json:"attemptedLanguages,omitempty"`
	Field              string   `json:"field,omitempty"`
}

func main() {
	js.Global().Set("__dataLocalizerInvoke", js.FuncOf(invoke))
	select {}
}

func invoke(_ js.Value, arguments []js.Value) any {
	if len(arguments) != 1 || arguments[0].Type() != js.TypeString {
		return encodeResponse(response{OK: false, Error: &responseError{
			Name: "TypeError", Message: "the Go bridge expects one JSON string",
		}})
	}

	var input request
	if err := json.Unmarshal([]byte(arguments[0].String()), &input); err != nil {
		return encodeError(err)
	}
	data, err := decodeJSON(input.Data)
	if err != nil {
		return encodeError(err)
	}

	var result any
	if input.Options == nil {
		result, err = localizer.Localize(data, input.LanguageHeader, input.FallbackLanguage)
	} else {
		policy, policyErr := parsePolicy(input.Options.MissingTranslation)
		if policyErr != nil {
			return encodeError(policyErr)
		}
		engine, createErr := localizer.New(localizer.Options{
			SupportedLanguages: input.Options.SupportedLanguages,
			FallbackLanguage:   input.Options.FallbackLanguage,
			FallbackLanguages:  input.Options.FallbackLanguages,
			MissingTranslation: policy,
		})
		if createErr != nil {
			return encodeError(createErr)
		}
		result, err = engine.Localize(data, input.LanguageHeader)
	}
	if err != nil {
		return encodeError(err)
	}
	return encodeResponse(response{OK: true, Value: result})
}

func decodeJSON(raw json.RawMessage) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	return value, nil
}

func parsePolicy(value string) (localizer.MissingTranslationPolicy, error) {
	switch value {
	case "", "preserve":
		return localizer.MissingPreserve, nil
	case "empty":
		return localizer.MissingEmpty, nil
	case "null":
		return localizer.MissingNull, nil
	case "error":
		return localizer.MissingError, nil
	default:
		return 0, &localizer.ConfigurationError{Field: "MissingTranslation", Message: "unknown policy " + value}
	}
}

func encodeError(err error) string {
	details := &responseError{Name: "Error", Message: err.Error()}
	var configuration *localizer.ConfigurationError
	var missing *localizer.MissingTranslationError
	var cycle *localizer.CycleError
	switch {
	case errors.As(err, &configuration):
		details.Name = "ConfigurationError"
		details.Field = configuration.Field
	case errors.As(err, &missing):
		details.Name = "MissingTranslationError"
		details.Path = missing.Path
		details.AttemptedLanguages = missing.AttemptedLanguages
	case errors.As(err, &cycle):
		details.Name = "CycleError"
		details.Path = cycle.Path
		details.OriginalPath = cycle.OriginalPath
	}
	return encodeResponse(response{OK: false, Error: details})
}

func encodeResponse(value response) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		fallback, _ := json.Marshal(response{OK: false, Error: &responseError{Name: "Error", Message: err.Error()}})
		return string(fallback)
	}
	return string(encoded)
}
