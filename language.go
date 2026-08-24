package localizer

import (
	"sort"
	"strconv"
	"strings"
)

type LanguagePreference struct {
	Tag     string
	Quality float64
	Order   int
}

func NormalizeLanguageTag(tag string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(tag), "_", "-"))
}

func validLanguageTag(tag string) bool {
	if tag == "*" {
		return true
	}
	parts := strings.Split(tag, "-")
	if len(parts) == 0 || len(parts[0]) < 1 || len(parts[0]) > 8 {
		return false
	}
	for index, part := range parts {
		if len(part) < 1 || len(part) > 8 {
			return false
		}
		for _, character := range part {
			isLetter := character >= 'a' && character <= 'z'
			isDigit := character >= '0' && character <= '9'
			if !isLetter && (index == 0 || !isDigit) {
				return false
			}
		}
	}
	return true
}

func ParseAcceptLanguage(header string) []LanguagePreference {
	if strings.TrimSpace(header) == "" {
		return nil
	}

	preferences := make([]LanguagePreference, 0)
	for order, rawEntry := range strings.Split(header, ",") {
		parts := strings.Split(rawEntry, ";")
		tag := NormalizeLanguageTag(parts[0])
		if !validLanguageTag(tag) {
			continue
		}

		quality := 1.0
		valid := true
		for _, parameter := range parts[1:] {
			parameter = strings.TrimSpace(parameter)
			name, rawValue, found := strings.Cut(parameter, "=")
			if !found || !strings.EqualFold(strings.TrimSpace(name), "q") {
				continue
			}
			value := strings.TrimSpace(rawValue)
			if !validQuality(value) {
				valid = false
				break
			}
			quality, _ = strconv.ParseFloat(value, 64)
		}
		if valid && quality > 0 {
			preferences = append(preferences, LanguagePreference{Tag: tag, Quality: quality, Order: order})
		}
	}

	sort.SliceStable(preferences, func(i, j int) bool {
		return preferences[i].Quality > preferences[j].Quality
	})
	return preferences
}

func validQuality(value string) bool {
	if value == "0" || value == "1" {
		return true
	}
	if strings.HasPrefix(value, "0.") {
		digits := value[2:]
		return len(digits) <= 3 && onlyDigits(digits)
	}
	if strings.HasPrefix(value, "1.") {
		digits := value[2:]
		return len(digits) <= 3 && onlyDigits(digits) && strings.Trim(digits, "0") == ""
	}
	return false
}

func onlyDigits(value string) bool {
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}

func negotiate(header string, supported map[string]struct{}, fallbacks []string) ([]string, string) {
	preferences := ParseAcceptLanguage(header)
	candidates := make([]string, 0, len(preferences)+len(fallbacks))
	requested := ""
	if len(preferences) > 0 {
		requested = preferences[0].Tag
	}

	appendUnique := func(language string) {
		if _, exists := supported[language]; !exists {
			return
		}
		for _, candidate := range candidates {
			if candidate == language {
				return
			}
		}
		candidates = append(candidates, language)
	}

	for _, preference := range preferences {
		if preference.Tag == "*" {
			continue
		}
		appendUnique(preference.Tag)
		if separator := strings.IndexByte(preference.Tag, '-'); separator > 0 {
			appendUnique(preference.Tag[:separator])
		}
	}
	for _, fallback := range fallbacks {
		appendUnique(fallback)
	}
	return candidates, requested
}
