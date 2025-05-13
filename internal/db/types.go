package db

type Place struct {
	GeonameID      int64   `json:"geonameid"`
	Name           string  `json:"name"`
	Asciiname      string  `json:"asciiname"`
	Alternatenames string  `json:"alternatenames"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	CountryCode    string  `json:"country_code"`
}
