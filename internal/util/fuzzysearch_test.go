package utils

import (
	"testing"

	"github.com/WhilCodingDoLeanr/iam_geocoder/internal/searchbarcache"
	"github.com/stretchr/testify/assert"
)

func TestFuzzySearch(t *testing.T) {
	err := searchbarcache.LoadCityIndex("city_config.json")
	assert.Nil(t, err)
	data := FuzzySearchByPrefix(searchbarcache.CityIndex, "at", 4)
	assert.NotNil(t, data)
}
