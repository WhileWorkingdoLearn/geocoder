package utils

import "github.com/WhilCodingDoLeanr/iam_geocoder/internal/config"

func BestToken(tokens []string) string {
	best := ""
	for _, t := range tokens {
		if config.Search.Stopwords[t] {
			continue
		}
		if config.Search.CountryCodes[t] {
			continue
		}
		if len(t) > len(best) {
			best = t
		}
	}
	// Fallback: erster Token
	if best == "" && len(tokens) > 0 {
		return tokens[0]
	}
	return best
}
