package models

import "time"

// The * means it's a pointer - can be nil (like NULL in SQL).
type BoardGame struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name" binding:"required"`
	MinPlayers    int       `json:"min_players" binding:"required"`
	MaxPlayers    int       `json:"max_players,omitempty"` // NULL in DB
	PlayTime      int       `json:"play_time" binding:"required"`
	MinAge        int       `json:"min_age" binding:"required"`
	Description   string    `json:"description" binding:"required"`
	ImageIDs      []int64   `json:"image_ids,omitempty"`
	CoverImageUrL string    `json:"coverImageUrl,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BoardGameImage struct {
	ID            int64
	BoardGameID   int64
	ImageData     []byte
	ImageMimeType string
	ThumbnailData []byte
	ImageType     string
	DisplayOrder  int
	UploadedAt    time.Time
}
