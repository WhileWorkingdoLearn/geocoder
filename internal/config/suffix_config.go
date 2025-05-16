package config

import (
	"encoding/json"
	"os"
	"strings"
)

var SuffixCfg *Config

func init() {
	SuffixCfg = &Config{}
}

func LoadSuffixesFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	cfg := SuffixCfg
	err = json.Unmarshal(data, &cfg.suffixes)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) GetSuffixesForLanguage(lang string) []string {
	lang = strings.ToLower(lang)
	if sfx, ok := c.suffixes[lang]; ok {
		return sfx
	}
	return nil
}
