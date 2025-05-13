-- Datei: init/01_create_places.sql

CREATE TABLE geonames_reduced (
    geonameid      BIGINT PRIMARY KEY,
    asciiname      TEXT,
    latitude       DOUBLE PRECISION,
    longitude      DOUBLE PRECISION,
    country_code   CHAR(2),
    admin1_code    TEXT,
    admin2_code    TEXT,
    admin3_code    TEXT,
    admin4_code    TEXT,
    timezone       TEXT,
    token          TEXT,
    location       GEOGRAPHY(Point, 4326) GENERATED ALWAYS AS (
        ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)::geography
    ) STORED
);


CREATE INDEX idx_geonames_token ON geonames_reduced USING GIN (token gin_trgm_ops);
CREATE INDEX idx_geonames_location ON geonames_reduced USING GIST (location);
