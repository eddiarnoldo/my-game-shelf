package models

import "time"

// The * means it's a pointer - can be nil (like NULL in SQL).
type BoardGame struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	MinPlayers  int       `json:"min_players"`
	MaxPlayers  *int      `json:"max_players,omitempty"` // NULL in DB
	PlayTime    *int      `json:"play_time,omitempty"`   // NULL in DB
	MinAge      int       `json:"min_age"`
	Description *string   `json:"description,omitempty"` // NULL in DB
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
