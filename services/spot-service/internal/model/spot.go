package model

type Spot struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Lat          float64  `json:"lat"`
	Lon          float64  `json:"lon"`
	Type         string   `json:"type"`
	AccessNotes  string   `json:"access_notes"`
	Difficulty   string   `json:"difficulty"`
	Rating       float64  `json:"rating"`
	ReviewCount  int      `json:"review_count"`
	Species      []string `json:"species"`
	BestSeasons  []string `json:"best_seasons"`
	CreatedAt    string   `json:"created_at"`
}

type SpotDetail struct {
	Spot
	Description   string   `json:"description"`
	DepthInfo     string   `json:"depth_info"`
	BottomType    string   `json:"bottom_type"`
	Facilities    []string `json:"facilities"`
	Regulations   string   `json:"regulations"`
	Parking       string   `json:"parking"`
	Lat           float64  `json:"lat"`
	Lon           float64  `json:"lon"`
}

type ListSpotsParams struct {
	Type      string
	Species   string
	Difficulty string
	Page      int
	Limit     int
}

type NearbyParams struct {
	Lat      float64
	Lon      float64
	RadiusMi float64
	Limit    int
}

type CreateSpotRequest struct {
	Name       string  `json:"name"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	Type       string  `json:"type"`
	Difficulty string  `json:"difficulty"`
}
