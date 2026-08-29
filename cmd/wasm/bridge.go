//go:build js && wasm

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"syscall/js"

	localizer "github.com/elrefai99/data-localizer"
)

type bridge struct {
	functions               []js.Func
	dataLocalizerError      js.Value
	configurationError      js.Value
	missingTranslationError js.Value
	cycleError              js.Value
}

type requestLocalizer struct {
	bridge             *bridge
	engine             *localizer.Localizer
	readLanguage       js.Value
	localizeForRequest js.Value
}

func newBridge() *bridge {
	b := &bridge{}
	b.dataLocalizerError = b.errorClass("DataLocalizerError", js.Global().Get("Error"))
	b.configurationError = b.errorClass("ConfigurationError", b.dataLocalizerError)
	b.missingTranslationError = b.errorClass("MissingTranslationError", b.dataLocalizerError)
	b.cycleError = b.errorClass("CycleError", b.dataLocalizerError)
	return b
}

func (b *bridge) modules() js.Value {
	modules := object()
	modules.Set("core", b.coreModule())
	modules.Set("framework", b.frameworkModule())
	modules.Set("express", b.expressModule())
	modules.Set("nest", b.nestModule())
	modules.Set("fastify", b.fastifyModule())
	modules.Set("koa", b.koaModule())
	return modules
}

func (b *bridge) coreModule() js.Value {
	module := object()
	module.Set("localize", b.function(b.localize))
	module.Set("createLocalizer", b.function(b.createLocalizer))
	module.Set("DataLocalizerError", b.dataLocalizerError)
	module.Set("ConfigurationError", b.configurationError)
	module.Set("MissingTranslationError", b.missingTranslationError)
	module.Set("CycleError", b.cycleError)
	return module
}

func (b *bridge) localize(_ js.Value, arguments []js.Value) any {
	return b.safely(func() any {
		if len(arguments) == 0 {
			return b.typeError("data must be JSON-compatible")
		}
		data, err := decodeValue(arguments[0])
		if err != nil {
			return b.typeError("data must be JSON-compatible: " + err.Error())
		}
		header := stringArgument(arguments, 1, "")
		fallback := stringArgument(arguments, 2, "en")
		result, err := localizer.Localize(data, header, fallback)
		return b.result(result, err)
	})
}

func (b *bridge) createLocalizer(_ js.Value, arguments []js.Value) any {
	return b.safely(func() any {
		engine, err := newLocalizer(valueArgument(arguments, 0))
		if err != nil {
			return b.errorValue(err)
		}
		result := object()
		result.Set("localize", b.function(func(_ js.Value, arguments []js.Value) any {
			return b.localizeWith(engine, arguments)
		}))
		return js.Global().Get("Object").Call("freeze", result)
	})
}

func (b *bridge) localizeWith(engine *localizer.Localizer, arguments []js.Value) any {
	return b.safely(func() any {
		if len(arguments) == 0 {
			return b.typeError("data must be JSON-compatible")
		}
		data, err := decodeValue(arguments[0])
		if err != nil {
			return b.typeError("data must be JSON-compatible: " + err.Error())
		}
		result, err := engine.Localize(data, stringArgument(arguments, 1, ""))
		return b.result(result, err)
	})
}

func newLocalizer(options js.Value) (*localizer.Localizer, error) {
	if isNullish(options) {
		return localizer.New(localizer.Options{})
	}
	if options.Type() != js.TypeObject {
		return nil, &localizer.ConfigurationError{Field: "Options", Message: "must be an object"}
	}
	supported, err := stringSliceProperty(options, "supportedLanguages")
	if err != nil {
		return nil, err
	}
	fallbacks, err := stringSliceProperty(options, "fallbackLanguages")
	if err != nil {
		return nil, err
	}
	policy, err := parsePolicy(stringProperty(options, "missingTranslation"))
	if err != nil {
		return nil, err
	}
	return localizer.New(localizer.Options{
		SupportedLanguages: supported,
		FallbackLanguage:   stringProperty(options, "fallbackLanguage"),
		FallbackLanguages:  fallbacks,
		MissingTranslation: policy,
	})
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

func (b *bridge) newRequestLocalizer(options js.Value) (*requestLocalizer, error) {
	engine, err := newLocalizer(options)
	if err != nil {
		return nil, err
	}
	readLanguage := js.Undefined()
	if !isNullish(options) {
		readLanguage = options.Get("getLanguageHeader")
		if !isNullish(readLanguage) && readLanguage.Type() != js.TypeFunction {
			return nil, &localizer.ConfigurationError{Field: "getLanguageHeader", Message: "must be a function"}
		}
	}
	result := &requestLocalizer{bridge: b, engine: engine, readLanguage: readLanguage}
	result.localizeForRequest = b.function(func(this js.Value, arguments []js.Value) any {
		return result.localizeWithRequest(arguments, this)
	})
	return result, nil
}

func (r *requestLocalizer) languageFor(request js.Value) string {
	if r.readLanguage.Type() == js.TypeFunction {
		value := r.readLanguage.Invoke(request)
		if value.Type() == js.TypeString {
			return value.String()
		}
		return ""
	}
	return acceptLanguage(request)
}

func (r *requestLocalizer) localize(arguments []js.Value) any {
	if len(arguments) == 0 {
		return r.bridge.typeError("data must be JSON-compatible")
	}
	header := ""
	if len(arguments) > 1 {
		if arguments[1].Type() == js.TypeString {
			header = arguments[1].String()
		} else if !isNullish(arguments[1]) {
			header = r.languageFor(arguments[1])
		}
	}
	return r.bridge.localizeWith(r.engine, []js.Value{arguments[0], js.ValueOf(header)})
}

func (r *requestLocalizer) forRequest(request js.Value) js.Value {
	return r.localizeForRequest.Call("bind", request)
}

func acceptLanguage(request js.Value) string {
	if request.Type() == js.TypeString {
		return request.String()
	}
	if isNullish(request) || request.Type() != js.TypeObject {
		return ""
	}
	if getter := request.Get("get"); getter.Type() == js.TypeFunction {
		if value := request.Call("get", "Accept-Language"); value.Type() == js.TypeString {
			return value.String()
		}
	}
	headers := request.Get("headers")
	if isNullish(headers) || headers.Type() != js.TypeObject {
		return ""
	}
	if getter := headers.Get("get"); getter.Type() == js.TypeFunction {
		value := headers.Call("get", "accept-language")
		if value.Type() == js.TypeString {
			return value.String()
		}
		return ""
	}
	keys := js.Global().Get("Object").Call("keys", headers)
	for index := 0; index < keys.Length(); index++ {
		key := keys.Index(index).String()
		if !strings.EqualFold(key, "accept-language") {
			continue
		}
		value := headers.Get(key)
		if js.Global().Get("Array").Call("isArray", value).Bool() {
			return value.Call("join", ",").String()
		}
		if value.Type() == js.TypeString {
			return value.String()
		}
	}
	return ""
}

func decodeValue(value js.Value) (_ any, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("%v", recovered)
		}
	}()
	encoded := js.Global().Get("JSON").Call("stringify", value)
	if encoded.Type() != js.TypeString {
		return nil, errors.New("JSON.stringify returned undefined")
	}
	decoder := json.NewDecoder(bytes.NewBufferString(encoded.String()))
	decoder.UseNumber()
	var result any
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func encodeValue(value any) (js.Value, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return js.Undefined(), err
	}
	return js.Global().Get("JSON").Call("parse", string(encoded)), nil
}

func (b *bridge) result(value any, err error) any {
	if err != nil {
		return b.errorValue(err)
	}
	encoded, err := encodeValue(value)
	if err != nil {
		return b.errorValue(err)
	}
	return encoded
}

func (b *bridge) errorValue(err error) js.Value {
	constructor := b.dataLocalizerError
	var configuration *localizer.ConfigurationError
	var missing *localizer.MissingTranslationError
	var cycle *localizer.CycleError
	switch {
	case errors.As(err, &configuration):
		constructor = b.configurationError
	case errors.As(err, &missing):
		constructor = b.missingTranslationError
	case errors.As(err, &cycle):
		constructor = b.cycleError
	}
	value := constructor.New(err.Error())
	switch {
	case configuration != nil:
		value.Set("field", configuration.Field)
	case missing != nil:
		value.Set("path", missing.Path)
		value.Set("attemptedLanguages", stringsValue(missing.AttemptedLanguages))
	case cycle != nil:
		value.Set("path", cycle.Path)
		value.Set("originalPath", cycle.OriginalPath)
	}
	return value
}

func (b *bridge) typeError(message string) js.Value {
	return js.Global().Get("TypeError").New(message)
}

func (b *bridge) errorClass(name string, parent js.Value) js.Value {
	var constructor js.Value
	constructor = b.function(func(_ js.Value, arguments []js.Value) any {
		value := js.Global().Get("Error").New(stringArgument(arguments, 0, ""))
		js.Global().Get("Object").Call("setPrototypeOf", value, constructor.Get("prototype"))
		value.Set("name", name)
		return value
	})
	prototype := js.Global().Get("Object").Call("create", parent.Get("prototype"))
	constructor.Set("prototype", prototype)
	js.Global().Get("Object").Call("defineProperty", prototype, "constructor", map[string]any{
		"value": constructor, "writable": true, "configurable": true,
	})
	return constructor
}

func (b *bridge) function(callback func(js.Value, []js.Value) any) js.Value {
	function := js.FuncOf(callback)
	b.functions = append(b.functions, function)
	return function.Value
}

func (b *bridge) wrapped(function js.Value) js.Value {
	return js.Global().Get("__dataLocalizerWrap").Invoke(function)
}

func (b *bridge) safely(callback func() any) (result any) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = b.typeError(fmt.Sprint(recovered))
		}
	}()
	return callback()
}

func object() js.Value {
	return js.Global().Get("Object").New()
}

func valueArgument(arguments []js.Value, index int) js.Value {
	if index >= len(arguments) {
		return js.Undefined()
	}
	return arguments[index]
}

func stringArgument(arguments []js.Value, index int, fallback string) string {
	if index >= len(arguments) || isNullish(arguments[index]) || arguments[index].Type() != js.TypeString {
		return fallback
	}
	return arguments[index].String()
}

func stringProperty(value js.Value, name string) string {
	property := value.Get(name)
	if property.Type() == js.TypeString {
		return property.String()
	}
	return ""
}

func stringSliceProperty(value js.Value, name string) ([]string, error) {
	property := value.Get(name)
	if isNullish(property) {
		return nil, nil
	}
	if !js.Global().Get("Array").Call("isArray", property).Bool() {
		return nil, &localizer.ConfigurationError{Field: name, Message: "must be an array of strings"}
	}
	result := make([]string, property.Length())
	for index := range result {
		if property.Index(index).Type() != js.TypeString {
			return nil, &localizer.ConfigurationError{Field: name, Message: "must be an array of strings"}
		}
		result[index] = property.Index(index).String()
	}
	return result, nil
}

func stringsValue(values []string) js.Value {
	result := js.Global().Get("Array").New(len(values))
	for index, value := range values {
		result.SetIndex(index, value)
	}
	return result
}

func isNullish(value js.Value) bool {
	return value.Type() == js.TypeUndefined || value.Type() == js.TypeNull
}
