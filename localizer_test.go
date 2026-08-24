package localizer

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestLocalizeNestedDataWithoutMutation(t *testing.T) {
	input := map[string]any{
		"title": map[string]any{"ar": "عنوان", "en": "Title"},
		"posts": []any{
			map[string]any{"name": map[string]any{"ar": "الأول", "en": "First"}, "published": true},
			map[string]any{"name": map[string]any{"ar": "الثاني", "en": "Second"}, "published": false},
		},
		"count": 2,
	}
	want := map[string]any{
		"title": "عنوان",
		"posts": []any{
			map[string]any{"name": "الأول", "published": true},
			map[string]any{"name": "الثاني", "published": false},
		},
		"count": 2,
	}

	got, err := Localize(input, " ar-EG, en;q=0.8 ")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Localize() = %#v, want %#v", got, want)
	}
	if input["title"].(map[string]any)["ar"] != "عنوان" {
		t.Fatal("Localize mutated its input")
	}
	gotMap := got.(map[string]any)
	if reflect.ValueOf(gotMap).Pointer() == reflect.ValueOf(input).Pointer() {
		t.Fatal("Localize returned the original map")
	}
}

func TestLocalizeUsesQualityValues(t *testing.T) {
	input := map[string]any{"greeting": map[string]any{"en": "Hello", "fr": "Bonjour"}}
	got, err := Localize(input, "en;q=0.5,fr;q=0.9")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"greeting": "Bonjour"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Localize() = %#v, want %#v", got, want)
	}
}

func TestTranslationsPreserveFalsyValues(t *testing.T) {
	engine := MustNew(Options{SupportedLanguages: []string{"en", "ar"}, FallbackLanguage: "en"})
	tests := []struct {
		name string
		in   any
		want any
	}{
		{name: "empty string", in: map[string]any{"ar": "", "en": "fallback"}, want: ""},
		{name: "false", in: map[string]any{"ar": false, "en": true}, want: false},
		{name: "zero", in: map[string]any{"ar": 0, "en": 1}, want: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := engine.Localize(test.in, "ar")
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("Localize() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestDefaultDetectorDoesNotCollapseDomainObjects(t *testing.T) {
	input := map[string]any{
		"routes": map[string]any{
			"en": map[string]any{"path": "/english", "enabled": true},
			"ar": map[string]any{"path": "/arabic", "enabled": false},
		},
	}
	got, err := Localize(input, "en")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, input) {
		t.Fatalf("Localize() = %#v, want ordinary object preserved", got)
	}
}

func TestMissingTranslationPolicies(t *testing.T) {
	input := map[string]any{"title": map[string]any{"ar": nil}}
	tests := []struct {
		policy MissingTranslationPolicy
		want   any
	}{
		{policy: MissingPreserve, want: map[string]any{"title": TranslationMap{"ar": nil}}},
		{policy: MissingEmpty, want: map[string]any{"title": ""}},
		{policy: MissingNull, want: map[string]any{"title": nil}},
	}
	for _, test := range tests {
		engine := MustNew(Options{
			SupportedLanguages: []string{"en", "ar"},
			FallbackLanguage:   "en",
			MissingTranslation: test.policy,
		})
		got, err := engine.Localize(input, "ar")
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("policy %v returned %#v, want %#v", test.policy, got, test.want)
		}
	}
}

func TestMissingErrorContainsPathAndLanguages(t *testing.T) {
	engine := MustNew(Options{
		SupportedLanguages: []string{"en", "ar"},
		FallbackLanguage:   "en",
		MissingTranslation: MissingError,
	})
	_, err := engine.Localize(map[string]any{"items": []any{map[string]any{"name": map[string]any{"ar": nil}}}}, "ar")
	var missing *MissingTranslationError
	if !errors.As(err, &missing) {
		t.Fatalf("error = %v, want MissingTranslationError", err)
	}
	if missing.Path != "$.items[0].name" || !reflect.DeepEqual(missing.AttemptedLanguages, []string{"ar", "en"}) {
		t.Fatalf("MissingTranslationError = %#v", missing)
	}
}

func TestCustomDetectorAndMissingHandler(t *testing.T) {
	engine := MustNew(Options{
		SupportedLanguages: []string{"en"},
		FallbackLanguage:   "en",
		IsTranslationMap: func(value TranslationMap, context Context) bool {
			return context.Path == "$.label"
		},
		OnMissing: func(value TranslationMap, context Context) (any, error) {
			return "missing at " + context.Path, nil
		},
	})
	got, err := engine.Localize(map[string]any{"label": map[string]any{"other": "value"}}, "en")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"label": "missing at $.label"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Localize() = %#v, want %#v", got, want)
	}
}

func TestTypedMapsSlicesAndSpecialValues(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	input := map[string]any{
		"labels": []map[string]string{{"en": "One", "ar": "واحد"}},
		"time":   now,
	}
	got, err := Localize(input, "en")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"labels": []any{"One"}, "time": now}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Localize() = %#v, want %#v", got, want)
	}
}

func TestCycleReturnsTypedError(t *testing.T) {
	input := map[string]any{}
	input["self"] = input
	_, err := Localize(input, "en")
	var cycle *CycleError
	if !errors.As(err, &cycle) {
		t.Fatalf("error = %v, want CycleError", err)
	}
	if cycle.Path != "$.self" || cycle.OriginalPath != "$" {
		t.Fatalf("CycleError = %#v", cycle)
	}
}

func TestNewRejectsUnsupportedFallback(t *testing.T) {
	_, err := New(Options{SupportedLanguages: []string{"ar"}, FallbackLanguage: "en"})
	var configuration *ConfigurationError
	if !errors.As(err, &configuration) {
		t.Fatalf("error = %v, want ConfigurationError", err)
	}
}
