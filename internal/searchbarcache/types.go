package searchbarcache

type City struct {
	Name        string     `json:"name"`
	Country     string     `json:"country"`
	Coordinates [2]float64 `json:"coordinates"`
}

type CityList struct {
	Cities []City `json:"cities"`
}
