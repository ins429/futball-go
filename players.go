package main

type PlayerNames struct {
	Names []string `form:"names" json:"names"`
}

type WCPlayerStatsResponse struct {
	Status int            `json:"status"`
	Stats  []WCPlayerStat `json:"stats"`
}

type WCPlayerStat struct {
	Name         string  `json:"name"`
	Age          float64 `json:"age"`
	Foot         string  `json:"foot"`
	BirthDate    string  `json:"birthDate"`
	BirthCountry string  `json:"birthCountry"`
	BirthCity    string  `json:"birthCity"`
	National     string  `json:"national"`
	Height       float64 `json:"height"`
	Weight       float64 `json:"weight"`
	Position     string  `json:"position"`
	Club         string  `json:"club"`
	Goals        float64 `json:"goals"`
	Assists      float64 `json:"assists"`
	PenaltyGoals float64 `json:"penaltyGoals"`
	OwnGoals     float64 `json:"ownGoals"`
	Image        string  `json:"image"`
}
