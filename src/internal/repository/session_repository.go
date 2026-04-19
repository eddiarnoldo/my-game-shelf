package repository

import (
	"context"
	"errors"

	"github.com/eddiarnoldo/my-game-shelf/src/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepo interface {
	Create(ctx context.Context, session *models.Session) error
	GetByID(ctx context.Context, id string) (*models.Session, error)
	Delete(ctx context.Context, id string) error
}

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	query := `INSERT INTO sessions (id, user_id, expires_at)
	          VALUES ($1, $2, $3)
	          RETURNING created_at`

	return r.db.QueryRow(ctx, query,
		session.ID,
		session.UserID,
		session.ExpiresAt,
	).Scan(&session.CreatedAt)
}

func (r *SessionRepository) GetByID(ctx context.Context, id string) (*models.Session, error) {
	query := `SELECT id, user_id, expires_at, created_at
	          FROM sessions WHERE id = $1`

	var s models.Session
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID,
		&s.UserID,
		&s.ExpiresAt,
		&s.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, id)
	return err
}
