package model

type SearchParams struct {
	Query    string
	Type     string
	Species  string
	State    string
	Difficulty string
	MinRating float64
	Page     int
	Limit    int
}

type SearchResult struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Lat         float64  `json:"lat"`
	Lon         float64  `json:"lon"`
	Species     []string `json:"species"`
	Rating      float64  `json:"rating"`
	Distance    *float64 `json:"distance,omitempty"`
}

type SpeciesResult struct {
	Name       string `json:"name"`
	CommonName string `json:"common_name"`
	SpotCount  int    `json:"spot_count"`
}
