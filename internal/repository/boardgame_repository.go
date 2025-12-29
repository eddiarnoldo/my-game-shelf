package repository

import (
	"context"

	"gihub.com/jackc/pgx/v5/pgxpool"
	"github.com/eddiarnoldo/my-game-shelf/internal/models"
)

// Create the type here to have a pointer to the db pool
type BoardGameRepository struct {
	db *pgxpool.Pool
}

func NewBoardGameRepository(db *pgxpool.Pool) *BoardGameRepository {
	return &BoardGameRepository{db: db}
}

func (r *BoardGameRepository) Create(ctx context.Context, game models.BoardGame) error {
	query := `INSERT into board_games ()`
}
