package utils

import (
	"fmt"
	"strings"

	"github.com/WhilCodingDoLeanr/iam_geocoder/internal/config"
)

func NewAddressParser(cfg config.SuffixConfig) *AddressParser {
	return &AddressParser{
		streetSuffixes: cfg,
		MaxNGram:       3,
	}
}

// Score-Bewertung
func (p *AddressParser) scoreStreetNGram(ngram []string, country string) (int, []string) {
	score := 0
	reasons := []string{}

	score += len(ngram)
	reasons = append(reasons, fmt.Sprintf("Length: +%d", len(ngram)))
	last := normalizeAndReduce(ngram[len(ngram)-1])
	for _, suffix := range p.streetSuffixes.GetSuffixesForLanguage(country) {
		if strings.HasSuffix(last, suffix) {
			score += 5
			reasons = append(reasons, fmt.Sprintf("Suffix match: +5 (%s)", suffix))
			break
		}
	}

	return score, reasons
}

func (p *AddressParser) DetectStreets(input string, country string) []Detection {
	tokens := strings.Fields(input)
	if len(tokens) == 0 {
		return nil
	}

	n := len(tokens)
	var results []Detection

	for size := p.MaxNGram; size >= 1; size-- {
		for i := 0; i <= n-size; i++ {
			ngram := tokens[i : i+size]
			fmt.Println(ngram)
			score, reason := p.scoreStreetNGram(ngram, country)
			if score >= 6 {
				results = append(results, Detection{
					Text:   strings.Join(ngram, " "),
					Start:  i,
					End:    i + size - 1,
					Score:  score,
					Reason: reason,
				})
			}
		}
	}
	return results
}
