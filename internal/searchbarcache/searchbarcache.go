package searchbarcache

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Index: prefix -> cities
type SearchCache map[string][]City

var CityIndex SearchCache

func init() {
	CityIndex = make(map[string][]City)
}

func LoadCityIndex(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("Fehler beim Öffnen der Datei: %w", err)
	}
	defer file.Close()

	var raw CityList
	if err := json.NewDecoder(file).Decode(&raw); err != nil {
		return fmt.Errorf("Fehler beim Parsen der JSON: %w", err)
	}

	for _, entry := range raw.Cities {
		if entry.Name == "" {
			continue
		}
		city := City{
			Name:        entry.Name,
			Coordinates: entry.Coordinates,
			Country:     entry.Country,
		}
		prefix := strings.ToLower(string([]rune(entry.Name)[0]))
		CityIndex[prefix] = append(CityIndex[prefix], city)
	}

	return nil
}
