package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	localizer "github.com/elrefai99/data-localizer"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "data-localizer:", err)
		os.Exit(1)
	}
}

func run(arguments []string, stdin io.Reader, stdout io.Writer) error {
	flags := flag.NewFlagSet("data-localizer", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	language := flags.String("lang", "", "Accept-Language value")
	fallback := flags.String("fallback", "en", "fallback language")
	supported := flags.String("supported", "", "comma-separated supported language tags")
	missing := flags.String("missing", "preserve", "preserve, empty, null, or error")
	pretty := flags.Bool("pretty", false, "indent JSON output")
	if err := flags.Parse(arguments); err != nil {
		return err
	}
	if flags.NArg() > 1 {
		return errors.New("usage: data-localizer [flags] [input.json]")
	}

	input := stdin
	if flags.NArg() == 1 && flags.Arg(0) != "-" {
		file, err := os.Open(flags.Arg(0))
		if err != nil {
			return err
		}
		defer file.Close()
		input = file
	}

	policy, err := parsePolicy(*missing)
	if err != nil {
		return err
	}
	options := localizer.Options{FallbackLanguage: *fallback, MissingTranslation: policy}
	if strings.TrimSpace(*supported) != "" {
		options.SupportedLanguages = splitList(*supported)
	}
	engine, err := localizer.New(options)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(input)
	decoder.UseNumber()
	var data any
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("decode input JSON: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return errors.New("input contains more than one JSON value")
		}
		return fmt.Errorf("decode input JSON: %w", err)
	}

	result, err := engine.Localize(data, *language)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(stdout)
	encoder.SetEscapeHTML(false)
	if *pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(result)
}

func parsePolicy(value string) (localizer.MissingTranslationPolicy, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "preserve":
		return localizer.MissingPreserve, nil
	case "empty":
		return localizer.MissingEmpty, nil
	case "null":
		return localizer.MissingNull, nil
	case "error":
		return localizer.MissingError, nil
	default:
		return 0, fmt.Errorf("unknown missing policy %q", value)
	}
}

func splitList(value string) []string {
	result := make([]string, 0)
	for _, item := range strings.Split(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			result = append(result, item)
		}
	}
	return result
}
