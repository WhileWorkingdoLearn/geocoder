package csvconverter

import (
	"regexp"
	"strings"
	"unicode"
)

func NewNormalizer() Normalizer {
	suffix := []string{
		"straße", "strasse", "str",
		"gasse", "allee", "weg",
	}
	return Normalizer{
		filter:         "aeiouäöüéèêáàâíóú",
		streetSuffixes: suffix,
		pattern:        `(?i)[\s\-]?(` + strings.Join(suffix, "|") + `)\b$`,
		suffixMap: map[string]string{
			"straße":  "_a",
			"strasse": "_a",
			"str":     "_a",
			"gasse":   "_b",
			"allee":   "_c",
			"weg":     "_d",
		},
	}
}

func (n *Normalizer) NormalizeTokenFromStreet(name string) string {
	txt := strings.ToLower(name)
	txt = strings.TrimSpace(txt)
	txt = n.removeStreetSuffix(txt)
	txt = n.removeVowels(txt)
	return txt
}

func (n *Normalizer) removeStreetSuffix(s string) string {

	re := regexp.MustCompile(n.pattern)

	return strings.TrimSpace(re.ReplaceAllString(s, ""))
}

func (n *Normalizer) ReplaceStreetSuffixWithMarker(s string) string {
	for suffix, marker := range n.suffixMap {
		pattern := `(?i)[\s\-]?` + suffix + `\b$`
		re := regexp.MustCompile(pattern)
		if re.MatchString(s) {
			return re.ReplaceAllString(s, marker)
		}
	}

	return s
}

func (n *Normalizer) removeVowels(input string) string {
	var result strings.Builder

	for _, r := range input {
		if !strings.ContainsRune(n.filter, r) && !strings.ContainsRune(n.filter, unicode.ToLower(r)) {
			result.WriteRune(r)
		}
	}

	return result.String()
}
