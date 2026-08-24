package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunLocalizesStandardInput(t *testing.T) {
	input := strings.NewReader(`{"title":{"en":"Hello","ar":"مرحبا"},"count":2}`)
	var output bytes.Buffer
	if err := run([]string{"-lang", "ar"}, input, &output); err != nil {
		t.Fatal(err)
	}
	want := `{"count":2,"title":"مرحبا"}` + "\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

func TestRunRejectsMultipleJSONValues(t *testing.T) {
	err := run(nil, strings.NewReader(`{} {}`), &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "more than one JSON value") {
		t.Fatalf("error = %v", err)
	}
}
