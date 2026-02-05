package models

import "time"

// The * means it's a pointer - can be nil (like NULL in SQL).
type BoardGame struct {
	ID              int64               `json:"id"`
	Name            string              `json:"name" binding:"required"`
	MinPlayers      int                 `json:"min_players" binding:"required"`
	MaxPlayers      int                 `json:"max_players,omitempty"` // NULL in DB
	PlayTime        int                 `json:"play_time" binding:"required"`
	MinAge          int                 `json:"min_age" binding:"required"`
	Description     string              `json:"description" binding:"required"`
	CoverImageUrl   string              `json:"coverImageUrl,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	BoardGameImages []BoardGameImageDto `json:"boardGameImages,omitempty"`
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

type BoardGameImageDto struct {
	ID           int64  `json:"id"`
	ImageUrl     string `json:"imageUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	ImageType    string `json:"imageType"`
	DisplayOrder int    `json:"displayOrder"`
}
