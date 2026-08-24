package localizer

import (
	"reflect"
	"testing"
)

func TestParseAcceptLanguageOrdersByQualityAndHeaderOrder(t *testing.T) {
	preferences := ParseAcceptLanguage("fr-CA;q=0.7, ar;q=0.9, en;q=0.9, *;q=0.1")
	want := []LanguagePreference{
		{Tag: "ar", Quality: 0.9, Order: 1},
		{Tag: "en", Quality: 0.9, Order: 2},
		{Tag: "fr-ca", Quality: 0.7, Order: 0},
		{Tag: "*", Quality: 0.1, Order: 3},
	}
	if !reflect.DeepEqual(preferences, want) {
		t.Fatalf("ParseAcceptLanguage() = %#v, want %#v", preferences, want)
	}
}

func TestParseAcceptLanguageIgnoresInvalidAndRejectedEntries(t *testing.T) {
	preferences := ParseAcceptLanguage("en;q=0, bad tag, ar;q=2, fr;q=0.75")
	want := []LanguagePreference{{Tag: "fr", Quality: 0.75, Order: 3}}
	if !reflect.DeepEqual(preferences, want) {
		t.Fatalf("ParseAcceptLanguage() = %#v, want %#v", preferences, want)
	}
}

func TestNormalizeLanguageTag(t *testing.T) {
	if got := NormalizeLanguageTag(" AR_eg "); got != "ar-eg" {
		t.Fatalf("NormalizeLanguageTag() = %q, want ar-eg", got)
	}
}

func TestNegotiateUsesExactBaseAndFallbackLanguages(t *testing.T) {
	supported := map[string]struct{}{"fr-ca": {}, "fr": {}, "en": {}, "ar": {}}
	got, requested := negotiate("fr-CA;q=0.9,ar;q=0.5", supported, []string{"en"})
	want := []string{"fr-ca", "fr", "ar", "en"}
	if requested != "fr-ca" || !reflect.DeepEqual(got, want) {
		t.Fatalf("negotiate() = %#v, %q; want %#v, fr-ca", got, requested, want)
	}
}
