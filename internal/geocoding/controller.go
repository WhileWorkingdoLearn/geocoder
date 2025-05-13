package geocoding

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/WhilCodingDoLeanr/iam_geocoder/internal/db"
	"github.com/gin-gonic/gin"
)

type GeoController struct {
	db db.IDBService
}

func NewGeoOCntroller(db db.IDBService) *GeoController {
	return &GeoController{db: db}
}

func (gs *GeoController) GetNearestPlaceHandler(ctx *gin.Context) {
	latStr := ctx.Query("lat")
	lonStr := ctx.Query("lon")

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lon, err2 := strconv.ParseFloat(lonStr, 64)
	if err1 != nil || err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Parameter: lat und lon müssen Floats sein"})
		return
	}

	place, err := gs.db.ReverseLocation(lat, lon)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, place)
}

const (
	maxTokens   = 5
	maxTokenLen = 30
)

// Bereinigt und splittet Eingaben auf Sonderzeichen wie +, _ oder -
func prepareTokens(input string) []string {
	// Sonderzeichen als Trenner behandeln
	re := regexp.MustCompile(`[^\p{L}\p{N}]+`) // alles außer Buchstaben und Zahlen
	clean := re.ReplaceAllString(input, " ")   // ersetzt z. B. _, +, , durch Leerzeichen
	clean = strings.ToLower(strings.TrimSpace(clean))

	tokens := strings.Fields(clean)

	var result []string
	for _, token := range tokens {
		if len(token) > 0 && len(token) <= maxTokenLen {
			result = append(result, token)
		}
		if len(result) >= maxTokens {
			break
		}
	}
	return result
}

func (gs *GeoController) GetSearchHandler(ctx *gin.Context) {
	query := ctx.Query("query")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query-Parameter ist erforderlich"})
		return
	}

	tokens := prepareTokens(query)
	results, err := gs.db.SearchForPlace(tokens)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, results)
}

func (gs *GeoController) PingHandler(ctx *gin.Context) {
	if err := gs.db.Ping(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "offline", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
