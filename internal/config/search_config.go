package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type SearchConfig struct {
	Stopwords    map[string]bool
	CountryCodes map[string]bool
}

type Stopwords struct {
	DE []string `json:"de"`
	EN []string `json:"en"`
	ES []string `json:"es"`
	FR []string `json:"fr"`
}

type RawCity struct {
	Name        string
	Coordinates [2]float64
}

type rawConfig struct {
	Stopwords    Stopwords               `json:"stopwords"`
	CountryCodes []string                `json:"country_codes"`
	KnownCities  []map[string][2]float64 `json:"known_cities"`
}

var Search *SearchConfig

func LoadSearchConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Fehler beim Öffnen der Such-Konfigurationsdatei: %w", err)
	}
	defer f.Close()

	var raw rawConfig
	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return fmt.Errorf("Fehler beim Parsen von JSON: %w", err)
	}

	set := make(map[string]bool)
	for _, val := range raw.Stopwords.DE {
		set[val] = true
	}
	for _, val := range raw.Stopwords.EN {
		set[val] = true
	}
	for _, val := range raw.Stopwords.ES {
		set[val] = true
	}
	for _, val := range raw.Stopwords.FR {
		set[val] = true
	}

	cityMap := make(map[string][2]float64)
	for _, entry := range raw.KnownCities {
		for name, coords := range entry {
			cityMap[name] = coords
		}
	}

	Search = &SearchConfig{
		Stopwords:    set,
		CountryCodes: sliceToSet(raw.CountryCodes),
	}

	return nil
}

func sliceToSet(slice []string) map[string]bool {
	set := make(map[string]bool)
	for _, val := range slice {
		set[val] = true
	}
	return set
}
