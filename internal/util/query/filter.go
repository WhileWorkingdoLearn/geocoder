package utils

import "strings"

func removeVowels(s string) string {
	vokale := "aeiouäöüAEIOUÄÖÜ"
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(vokale, r) {
			return -1
		}
		return r
	}, s)
}

func normalizeToken(s string) string {
	s = strings.ToLower(s)
	replacements := map[string]string{
		"ä": "ae", "ö": "oe", "ü": "ue", "ß": "ss",
	}
	for old, newVal := range replacements {
		s = strings.ReplaceAll(s, old, newVal)
	}
	return s
}

func normalizeAndReduce(s string) string {
	return removeVowels(normalizeToken(s))
}
