package repository

import (
	"github.com/eddiarnoldo/my-game-shelf/src/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BoardGameImageRepository struct {
	db *pgxpool.Pool
}

type BoardGameImageRepo interface {
	SaveImage(image *models.BoardGameImage) error
	GetAllImagesForBoardGame(boardGameId int64, imageType string) ([]*models.BoardGameImage, error)
	DeleteImage(id int64) error
}

func NewBoardGameImageRepository(db *pgxpool.Pool) *BoardGameImageRepository {
	return &BoardGameImageRepository{db: db}
}
