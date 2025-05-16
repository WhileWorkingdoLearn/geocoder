package config

type SuffixConfig interface {
	GetSuffixesForLanguage(lang string) []string
}

type StreetSuffixes map[string][]string

type Config struct {
	suffixes StreetSuffixes
}
