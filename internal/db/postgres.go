package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	utils "github.com/WhilCodingDoLeanr/iam_geocoder/internal/util"
)

type IDBService interface {
	ReverseLocation(lat, lon float64) (*Place, error)
	SearchForPlace(query []string) ([]Place, error)
	Ping() error
}

type DBService struct {
	db *sql.DB
}

func NewDBService(path string) (*DBService, error) {
	var err error
	DB, err := sql.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("error with sql.Open &w", err)
	}
	DB.SetMaxOpenConns(10)
	DB.SetConnMaxLifetime(5 * time.Minute)
	return &DBService{db: DB}, nil
}

func ToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func (dbs *DBService) ReverseLocation(lat, lon float64) (*Place, error) {
	query := `
		SELECT geonameid, name, asciiname, alternatenames, latitude, longitude, country_code
		FROM geonames
		ORDER BY location <-> ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
		LIMIT 1
	`

	row := dbs.db.QueryRowContext(context.Background(), query, lon, lat)

	var place Place
	var altName sql.NullString

	err := row.Scan(
		&place.GeonameID,
		&place.Name,
		&place.Asciiname,
		&altName,
		&place.Latitude,
		&place.Longitude,
		&place.CountryCode,
	)

	if err != nil {
		return nil, fmt.Errorf("DB-Fehler bei Reverse-Geocoding: %w", err)
	}
	place.Alternatenames = ToString(altName)

	return &place, nil
}

func (dbs *DBService) SearchForPlace(tokens []string) ([]Place, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("Keine Tokens zur Suche angegeben")
	}
	fmt.Println(tokens)
	base := `
		WITH matches AS (
			SELECT geonameid, asciiname, alternatenames, latitude, longitude, country_code
			FROM geonames
			WHERE %s
			LIMIT 200
		)
		SELECT *,
		       GREATEST(similarity(asciiname, $1), similarity(alternatenames, $1)) AS score
		FROM matches
		ORDER BY score DESC
		LIMIT 20;
	`

	var whereParts []string
	var args []interface{}

	scoringToken := utils.BestToken(tokens)

	for i, token := range tokens {
		whereParts = append(whereParts, fmt.Sprintf("(asciiname %% $%d OR alternatenames %% $%d)", i+2, i+2))
		args = append(args, token)
	}

	whereClause := strings.Join(whereParts, " OR ")
	query := fmt.Sprintf(base, whereClause)

	args = append([]interface{}{scoringToken}, args...)

	rows, err := dbs.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("DB-Fehler bei Tokensuche: %w", err)
	}
	defer rows.Close()

	var results []Place
	for rows.Next() {
		var p Place
		var alt sql.NullString
		var score float64

		err := rows.Scan(
			&p.GeonameID,
			&p.Asciiname,
			&alt,
			&p.Latitude,
			&p.Longitude,
			&p.CountryCode,
			&score,
		)
		if err != nil {
			return nil, fmt.Errorf("Scan-Fehler: %w", err)
		}
		p.Alternatenames = ToString(alt)
		results = append(results, p)
	}

	return results, nil
}

func (dbs *DBService) Ping() error {
	return dbs.db.Ping()
}
