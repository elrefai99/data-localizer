package localizer

import "fmt"

type MissingTranslationPolicy uint8

const (
	MissingPreserve MissingTranslationPolicy = iota
	MissingEmpty
	MissingNull
	MissingError
)

type TranslationMap map[string]any

type Context struct {
	Path               string
	RequestedLanguage  string
	AttemptedLanguages []string
}

type Detector func(value TranslationMap, context Context) bool

type MissingHandler func(value TranslationMap, context Context) (any, error)

type Options struct {
	SupportedLanguages []string
	FallbackLanguage   string
	FallbackLanguages  []string
	MissingTranslation MissingTranslationPolicy
	IsTranslationMap   Detector
	OnMissing          MissingHandler
}

type ConfigurationError struct {
	Field   string
	Message string
}

func (e *ConfigurationError) Error() string {
	return fmt.Sprintf("invalid localizer configuration %q: %s", e.Field, e.Message)
}

type MissingTranslationError struct {
	Path               string
	AttemptedLanguages []string
}

func (e *MissingTranslationError) Error() string {
	attempted := "none"
	if len(e.AttemptedLanguages) > 0 {
		attempted = join(e.AttemptedLanguages, ", ")
	}
	return fmt.Sprintf("missing translation at %s; attempted: %s", e.Path, attempted)
}

type CycleError struct {
	Path         string
	OriginalPath string
}

func (e *CycleError) Error() string {
	return fmt.Sprintf("cycle at %s refers to %s", e.Path, e.OriginalPath)
}

func join(values []string, separator string) string {
	if len(values) == 0 {
		return ""
	}
	result := values[0]
	for _, value := range values[1:] {
		result += separator + value
	}
	return result
}
