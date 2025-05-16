package utils

import (
	"fmt"
	"testing"

	"github.com/WhilCodingDoLeanr/iam_geocoder/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRemoveVowels(t *testing.T) {
	str := removeVowels("Häuser")
	assert.Equal(t, "Hsr", str)
	str = removeVowels("HäÄAser")
	assert.Equal(t, "Hsr", str)
}

func TestNormalizeToken(t *testing.T) {
	str := normalizeToken("Häuser")
	assert.Equal(t, "haeuser", str)
}

func TestNormalizeAndRedice(t *testing.T) {

	str := normalizeAndReduce("Häuser")
	assert.Equal(t, "hsr", str)

	str = normalizeAndReduce("Häus1r")
	assert.Equal(t, "hs1r", str)

	str = normalizeAndReduce("Häu s1r")
	assert.Equal(t, "h s1r", str)
}

func TestAnalyzer(t *testing.T) {

	err := config.LoadSuffixesFromFile("../../../config/street_suffix_config.json")
	assert.Nil(t, err)
	cfg := config.SuffixCfg
	adressParser := NewAddressParser(cfg)
	e1 := "lt mnzr strss 15 frnkfrt"
	scoresDe := adressParser.DetectStreets(e1, "de")
	fmt.Println(scoresDe)
	scoresFr := adressParser.DetectStreets(e1, "fr")
	fmt.Println(scoresFr)
	assert.Greater(t, len(scoresDe), 100)

}
