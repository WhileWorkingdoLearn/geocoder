package main

import (
	"fmt"
	"log"

	utils "github.com/WhilCodingDoLeanr/iam_geocoder/internal/util/query"
)

const streetcfg = "../config/street_suffixes_config.json"
const countrycfg = "../config/country_suffixes_config.json"

func main() {

	parser, err := utils.NewAddressParser("fr", streetcfg, countrycfg)
	if err != nil {
		log.Fatal(err)
	}
	detections := parser.ParseInput("Deutschland Rue de Leandro 13 12209")
	fmt.Println(detections)
}
