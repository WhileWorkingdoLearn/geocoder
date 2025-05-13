package main

import (
	"log"
	"net/http"
	"os"

	"github.com/WhilCodingDoLeanr/iam_geocoder/internal/db"
	"github.com/WhilCodingDoLeanr/iam_geocoder/internal/geocoding"
	"github.com/WhilCodingDoLeanr/iam_geocoder/internal/searchbarcache"
	utils "github.com/WhilCodingDoLeanr/iam_geocoder/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	dbConn := os.Getenv("DB_URL")
	database, err := db.NewDBService(dbConn)
	if err != nil {
		log.Fatalf("❌ DB-Verbindung fehlgeschlagen: %v", err)
	}

	pathToCfg := os.Getenv("CFG_CITYNAMES")
	if len(pathToCfg) <= len(".json") {
		log.Fatal("❌ no valid config provided")
	}

	if pathToCfg[len(pathToCfg)-len(".json"):] != ".json" {
		log.Fatal("❌ file must be json")
	}

	if err := searchbarcache.LoadCityIndex(pathToCfg); err != nil {
		log.Fatalf("❌ Fehler beim Laden der Such-Konfiguration: %v", err)
	}

	ctrl := geocoding.NewGeoOCntroller(database)

	router := gin.Default()

	router.GET("/search", ctrl.GetSearchHandler)

	router.GET("/searchbar", func(ctx *gin.Context) {
		query := ctx.Query("query")
		if query == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameter"})
			return
		}
		matches := utils.FuzzySearchByPrefix(searchbarcache.CityIndex, query, 2)
		ctx.JSON(http.StatusOK, matches)
	})

	router.GET("/reverse", ctrl.GetNearestPlaceHandler)

	router.GET("/ping/db", ctrl.PingHandler)

	router.GET("/tokens", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, searchbarcache.CityIndex)
	})

	log.Println("🚀 Server läuft auf http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("❌ Serverfehler: %v", err)
	}
}
