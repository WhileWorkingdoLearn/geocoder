package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	err := LoadSuffixesFromFile("../../config/street_suffix_config.json")
	assert.Nil(t, err)
	suffixes := SuffixCfg.GetSuffixesForLanguage("De")
	assert.NotNil(t, suffixes)
	assert.Greater(t, len(SuffixCfg.suffixes), 0)
	assert.Equal(t, "strss", suffixes[0])
}
