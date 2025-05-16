package main

import (
	"fmt"
	"strings"
)

func removeVowels(s string) string {
	vokale := "áaeéiouäöüAEIOUÄÖÜ"
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(vokale, r) {
			return 95
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

func removeDobuleUndersoce(s string) string {
	s = strings.ReplaceAll(s, "__", "_")
	return s
}

func normalizeAndReduce(s string) string {
	return removeDobuleUndersoce(removeVowels(normalizeToken(s)))
}

type Token struct {
	input      string
	normalized string
	value      int
}

type Result struct {
	Text   string
	Start  int
	End    int
	Score  int
	Reason string
}

type Detection struct {
	test   string
	score  int
	reason []string
}

const n_gram = 4

func main() {
	data := []string{
		"Autoroute",
		"Route",
		"Rue",
		"Boulevard",
		"Avenue",
		"Allée",
		"Impasse",
		"Chemin",
		"Voie rapide",
		"Passage",
		"Voie",
		"Chaussée",
		"Ruelle",
		"Cours",
		"Esplanade",
		"Sentier",
		"Traverse",
		"Promenade",
		"Quai",
		"Grande Rue",
		"Périphérique",
		"Bretelle",
		"Carrefour",
	}

	result := make([]string, 0)

	for _, d := range data {
		d = normalizeAndReduce(d)
		result = append(result, d)
	}
	/*
			for index, d := range result {
				result[index] = fmt.Sprintf("$%v", d)
			}*
		fmt.Println(result)
	*/
	suffixmap := map[string][]string{
		"fr": result,
	}

	splitInput := strings.Fields("Aba Rue de  Autnomus la")

	converted := make([]Token, len(splitInput))

	for i, d := range splitInput {

		t := Token{input: d}
		t.normalized = normalizeAndReduce(d)
		if i == 0 {
			t.value += 5
		}
		if i == len(splitInput)-1 {
			t.value += 5
		}
		converted[i] = t
	}

	n := len(converted)
	var results []Detection

	for size := n; size >= 1; size-- {
		for i := 0; i <= n-size; i++ {
			ngram := converted[i : i+size]
			score, reason := score(ngram, suffixmap)
			if score >= len(converted) {
				results = append(results, Detection{score: score, reason: reason})
			}
		}
	}

	fmt.Println(results)
}

func score(t []Token, suffixes map[string][]string) (int, []string) {
	score := 0
	score += len(t)
	first := t[0]
	reasons := []string{}
	for _, suffix := range suffixes["fr"] {
		if strings.HasSuffix(first.normalized, suffix) {
			score += len(t)
			reasons = append(reasons, fmt.Sprintf("val: %s Suffix match: +5 (%s)", len(t), suffix))
			break
		}
	}

	return score, reasons
}
