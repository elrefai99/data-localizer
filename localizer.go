package localizer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Localizer struct {
	supported map[string]struct{}
	fallbacks []string
	policy    MissingTranslationPolicy
	detector  Detector
	onMissing MissingHandler
}

func New(options Options) (*Localizer, error) {
	languages := options.SupportedLanguages
	if len(languages) == 0 {
		languages = defaultLanguages
	}

	supported := make(map[string]struct{}, len(languages))
	for _, language := range languages {
		normalized := NormalizeLanguageTag(language)
		if !validLanguageTag(normalized) || normalized == "*" {
			return nil, &ConfigurationError{Field: "SupportedLanguages", Message: fmt.Sprintf("invalid language tag %q", language)}
		}
		supported[normalized] = struct{}{}
	}
	if len(supported) == 0 {
		return nil, &ConfigurationError{Field: "SupportedLanguages", Message: "must contain at least one language"}
	}

	primaryFallback := NormalizeLanguageTag(options.FallbackLanguage)
	if primaryFallback == "" {
		if _, hasEnglish := supported["en"]; hasEnglish {
			primaryFallback = "en"
		} else {
			normalized := make([]string, 0, len(supported))
			for language := range supported {
				normalized = append(normalized, language)
			}
			sort.Strings(normalized)
			primaryFallback = normalized[0]
		}
	}

	fallbacks := make([]string, 0, 1+len(options.FallbackLanguages))
	for _, fallback := range append([]string{primaryFallback}, options.FallbackLanguages...) {
		normalized := NormalizeLanguageTag(fallback)
		if _, exists := supported[normalized]; !exists {
			return nil, &ConfigurationError{Field: "FallbackLanguages", Message: fmt.Sprintf("unsupported language %q", fallback)}
		}
		if !contains(fallbacks, normalized) {
			fallbacks = append(fallbacks, normalized)
		}
	}

	if options.MissingTranslation > MissingError {
		return nil, &ConfigurationError{Field: "MissingTranslation", Message: "unknown policy"}
	}

	return &Localizer{
		supported: supported,
		fallbacks: fallbacks,
		policy:    options.MissingTranslation,
		detector:  options.IsTranslationMap,
		onMissing: options.OnMissing,
	}, nil
}

func MustNew(options Options) *Localizer {
	localizer, err := New(options)
	if err != nil {
		panic(err)
	}
	return localizer
}

func Localize(data any, languageHeader string, fallbackLanguage ...string) (any, error) {
	fallback := "en"
	if len(fallbackLanguage) > 0 && strings.TrimSpace(fallbackLanguage[0]) != "" {
		fallback = fallbackLanguage[0]
	}
	engine, err := New(Options{FallbackLanguage: fallback, MissingTranslation: MissingEmpty})
	if err != nil {
		return nil, err
	}
	return engine.Localize(data, languageHeader)
}

func (l *Localizer) Localize(data any, languageHeader string) (any, error) {
	candidates, requested := negotiate(languageHeader, l.supported, l.fallbacks)
	state := walkState{active: make(map[visit]string)}
	return l.walk(reflect.ValueOf(data), "$", requested, candidates, &state)
}

type visit struct {
	kind reflect.Kind
	ptr  uintptr
}

type walkState struct {
	active map[visit]string
}

func (l *Localizer) walk(value reflect.Value, path, requested string, candidates []string, state *walkState) (any, error) {
	if !value.IsValid() {
		return nil, nil
	}
	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return nil, nil
		}
		return l.walk(value.Elem(), path, requested, candidates, state)
	}

	switch value.Kind() {
	case reflect.Map:
		if value.Type().Key().Kind() != reflect.String {
			return value.Interface(), nil
		}
		if value.IsNil() {
			return nil, nil
		}
		key := visit{kind: reflect.Map, ptr: value.Pointer()}
		if original, exists := state.active[key]; exists && key.ptr != 0 {
			return nil, &CycleError{Path: path, OriginalPath: original}
		}
		if key.ptr != 0 {
			state.active[key] = path
			defer delete(state.active, key)
		}

		object := make(TranslationMap, value.Len())
		iterator := value.MapRange()
		for iterator.Next() {
			object[iterator.Key().String()] = iterator.Value().Interface()
		}
		context := Context{Path: path, RequestedLanguage: requested, AttemptedLanguages: cloneStrings(candidates)}
		if l.isTranslationMap(object, context) {
			return l.resolve(object, context)
		}

		result := make(map[string]any, len(object))
		for key, child := range object {
			localized, err := l.walk(reflect.ValueOf(child), childPath(path, key), requested, candidates, state)
			if err != nil {
				return nil, err
			}
			result[key] = localized
		}
		return result, nil

	case reflect.Slice, reflect.Array:
		if value.Kind() == reflect.Slice && value.IsNil() {
			return nil, nil
		}
		var key visit
		if value.Kind() == reflect.Slice {
			key = visit{kind: reflect.Slice, ptr: value.Pointer()}
			if original, exists := state.active[key]; exists && key.ptr != 0 {
				return nil, &CycleError{Path: path, OriginalPath: original}
			}
			if key.ptr != 0 {
				state.active[key] = path
				defer delete(state.active, key)
			}
		}
		result := make([]any, value.Len())
		for index := 0; index < value.Len(); index++ {
			localized, err := l.walk(value.Index(index), fmt.Sprintf("%s[%d]", path, index), requested, candidates, state)
			if err != nil {
				return nil, err
			}
			result[index] = localized
		}
		return result, nil

	case reflect.Pointer, reflect.Func, reflect.Chan, reflect.Struct:
		return value.Interface(), nil
	default:
		return value.Interface(), nil
	}
}

func (l *Localizer) isTranslationMap(value TranslationMap, context Context) bool {
	if l.detector != nil {
		return l.detector(cloneMap(value), context)
	}
	if len(value) == 0 {
		return false
	}
	seen := make(map[string]struct{}, len(value))
	for language, translation := range value {
		normalized := NormalizeLanguageTag(language)
		if _, exists := l.supported[normalized]; !exists || !isTranslationLeaf(translation) {
			return false
		}
		if _, duplicate := seen[normalized]; duplicate {
			return false
		}
		seen[normalized] = struct{}{}
	}
	return true
}

func isTranslationLeaf(value any) bool {
	if value == nil {
		return true
	}
	switch value.(type) {
	case string, bool, json.Number,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	default:
		return false
	}
}

func (l *Localizer) resolve(value TranslationMap, context Context) (any, error) {
	entries := make(map[string]any, len(value))
	for language, translation := range value {
		entries[NormalizeLanguageTag(language)] = translation
	}
	for _, language := range context.AttemptedLanguages {
		if translation, exists := entries[language]; exists && translation != nil {
			return translation, nil
		}
	}
	if l.onMissing != nil {
		return l.onMissing(cloneMap(value), context)
	}
	switch l.policy {
	case MissingPreserve:
		return cloneMap(value), nil
	case MissingEmpty:
		return "", nil
	case MissingNull:
		return nil, nil
	case MissingError:
		return nil, &MissingTranslationError{Path: context.Path, AttemptedLanguages: cloneStrings(context.AttemptedLanguages)}
	default:
		panic("unreachable missing-translation policy")
	}
}

func cloneMap(value TranslationMap) TranslationMap {
	result := make(TranslationMap, len(value))
	for key, child := range value {
		result[key] = child
	}
	return result
}

func cloneStrings(values []string) []string {
	return append([]string(nil), values...)
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func childPath(parent, key string) string {
	if validIdentifier(key) {
		return parent + "." + key
	}
	return parent + "[" + strconv.Quote(key) + "]"
}

func validIdentifier(value string) bool {
	if value == "" {
		return false
	}
	for index, character := range value {
		letter := character >= 'a' && character <= 'z' || character >= 'A' && character <= 'Z'
		if !letter && character != '_' && character != '$' && (index == 0 || character < '0' || character > '9') {
			return false
		}
	}
	return true
}
