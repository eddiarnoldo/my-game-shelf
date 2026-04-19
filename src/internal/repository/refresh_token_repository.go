package repository

import (
	"context"
	"errors"

	"github.com/eddiarnoldo/my-game-shelf/src/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepo interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	GetByTokenHash(ctx context.Context, hash string) (*models.RefreshToken, error)
	MarkUsed(ctx context.Context, id int64) error
	DeleteBySessionID(ctx context.Context, sessionID string) error
}

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (session_id, token_hash, expires_at)
	          VALUES ($1, $2, $3)
	          RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		token.SessionID,
		token.TokenHash,
		token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
}

func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, hash string) (*models.RefreshToken, error) {
	query := `SELECT id, session_id, token_hash, used, created_at, expires_at
	          FROM refresh_tokens WHERE token_hash = $1`

	var t models.RefreshToken
	err := r.db.QueryRow(ctx, query, hash).Scan(
		&t.ID,
		&t.SessionID,
		&t.TokenHash,
		&t.Used,
		&t.CreatedAt,
		&t.ExpiresAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTokenNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *RefreshTokenRepository) MarkUsed(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `UPDATE refresh_tokens SET used = TRUE WHERE id = $1`, id)
	return err
}

func (r *RefreshTokenRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE session_id = $1`, sessionID)
	return err
}
