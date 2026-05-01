package model

type Review struct {
	ID        string  `json:"id"`
	SpotID    string  `json:"spot_id"`
	UserID    string  `json:"user_id"`
	Username  string  `json:"username"`
	Rating    int     `json:"rating"`
	Content   string  `json:"content"`
	CreatedAt string  `json:"created_at"`
}

type CreateReviewRequest struct {
	Rating  int    `json:"rating"`
	Content string `json:"content"`
}

type CatchLog struct {
	ID        string   `json:"id"`
	UserID    string   `json:"user_id"`
	SpotID    string   `json:"spot_id"`
	SpotName  string   `json:"spot_name"`
	Species   string   `json:"species"`
	Weight    float64  `json:"weight_lbs"`
	Length    float64  `json:"length_in"`
	BaitUsed  string   `json:"bait_used"`
	Photos    []string `json:"photos"`
	CreatedAt string   `json:"created_at"`
}

type CreateCatchRequest struct {
	SpotID  string  `json:"spot_id"`
	Species string  `json:"species"`
	Weight  float64 `json:"weight_lbs"`
	Length  float64 `json:"length_in"`
	BaitUsed string `json:"bait_used"`
}
