package models

type BoardGame struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	MinPlayers  int    `json:"min_players"`
	MaxPlayers  int    `json:"max_players"`
	PlayTime    int    `json:"play_time"` // in minutes
	AgeRating   int    `json:"age_rating"`
	Description string `json:"description"`
}
