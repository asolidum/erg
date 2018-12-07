package models

type OpenCage struct {
	Status struct {
		Code int `json:"code"`
	} `json:"status"`
	Results []struct {
		Name string `json:"formatted"`
	} `json:"results"`
}

type FourSquare struct {
	Meta struct {
		Code int `json:"code"`
	} `json:"meta"`
	Response struct {
		Venues []struct {
			Name string `json:"name"`
		} `json:"venues"`
	} `json:"response"`
}

type Yelp struct {
	Businesses []struct {
		Name string `json:"name"`
	} `json:"businesses"`
}

type MapBox struct {
	Features []struct {
		Name string `json:"text"`
	}
}

type OSM struct {
	Name string `json:"name"`
}

type RGStats struct {
	Num          int
	Latitude     float64
	Longitude    float64
	RadiusMeters int // Not all rev geocoders support this
	OpenCage     string
	Fourquare    string
	Yelp         string
	MapBox       string
	OSM          string
}
