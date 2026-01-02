package repository

import (
	"context"

	"github.com/eddiarnoldo/my-game-shelf/src/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Create the type here to have a pointer to the db pool
type BoardGameRepository struct {
	db *pgxpool.Pool
}

func NewBoardGameRepository(db *pgxpool.Pool) *BoardGameRepository {
	return &BoardGameRepository{db: db}
}

func (r *BoardGameRepository) Create(ctx context.Context, game *models.BoardGame) error {
	query := `INSERT into board_games 
		(name, min_players, max_players, play_time, min_age, description)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`

	//Here we execute the query and assign the returned id and created_at to the game struct
	err := r.db.QueryRow(ctx, query,
		game.Name,
		game.MinPlayers,
		game.MaxPlayers,
		game.PlayTime,
		game.MinAge,
		game.Description,
	).Scan(&game.ID, &game.CreatedAt)

	return err
}

func (r *BoardGameRepository) GetAll(ctx context.Context) ([]*models.BoardGame, error) {
	query := `SELECT id, name, min_players, max_players, play_time, min_age, description, created_at, updated_at
		FROM board_games ORDER BY id ASC`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, ErrQueryFailed
	}

	//Need to close resultset
	defer rows.Close()

	var boardGames []*models.BoardGame
	for rows.Next() {
		boardGame := &models.BoardGame{}
		err := rows.Scan(
			&boardGame.ID,
			&boardGame.Name,
			&boardGame.MinPlayers,
			&boardGame.MaxPlayers,
			&boardGame.PlayTime,
			&boardGame.MinAge,
			&boardGame.Description,
			&boardGame.CreatedAt,
			&boardGame.UpdatedAt,
		)

		if err != nil {
			return nil, ErrQueryFailed
		}
		boardGames = append(boardGames, boardGame)
	}

	return boardGames, nil
}

func (r *BoardGameRepository) GetByID(ctx context.Context, id int64) (*models.BoardGame, error) {
	query := `SELECT id, name, min_players, max_players, play_time, min_age, description, created_at, updated_at
		FROM board_games WHERE id = $1`

	var game models.BoardGame

	err := r.db.QueryRow(ctx, query, id).Scan(
		&game.ID,
		&game.Name,
		&game.MinPlayers,
		&game.MaxPlayers,
		&game.PlayTime,
		&game.MinAge,
		&game.Description,
		&game.CreatedAt,
		&game.UpdatedAt,
	)

	if err != nil {
		return nil, ErrQueryFailed
	}

	return &game, nil
}

// Delete removes a board game by ID
func (r *BoardGameRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM board_games WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return ErrQueryFailed
	}

	if commandTag.RowsAffected() == 0 {
		return ErrBoardGameNotFound
	}

	return nil
}
