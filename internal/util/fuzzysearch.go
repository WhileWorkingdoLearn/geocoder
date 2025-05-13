package utils

import (
	"fmt"
	"strings"

	"github.com/WhilCodingDoLeanr/iam_geocoder/internal/searchbarcache"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

func FuzzySearchByPrefix(cache searchbarcache.SearchCache, query string, maxDistance int) []searchbarcache.City {
	query = strings.ToLower(query)
	if len(query) == 0 {
		return nil
	}

	prefix := string([]rune(query)[0])
	candidates := cache[prefix]

	if len(query) == 1 {
		return cache[prefix]
	}

	var results []searchbarcache.City
	for _, city := range candidates {
		dist := levenshtein.DistanceForStrings([]rune(query), []rune(city.Name), levenshtein.DefaultOptions)
		if dist <= maxDistance {
			fmt.Println(dist)
			fmt.Println(city.Name)
			results = append(results, city)
		}
	}
	return results
}
