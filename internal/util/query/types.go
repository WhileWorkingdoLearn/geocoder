package utils

import "github.com/WhilCodingDoLeanr/iam_geocoder/internal/config"

type AddressParser struct {
	streetSuffixes config.SuffixConfig
	MaxNGram       int
}

type Detection struct {
	Text   string
	Start  int
	End    int
	Score  int
	Reason []string
}
