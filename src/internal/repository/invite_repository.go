package repository

import (
	"context"
	"errors"

	"github.com/eddiarnoldo/my-game-shelf/src/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InviteRepo interface {
	Create(ctx context.Context, invite *models.Invite) error
	GetByCode(ctx context.Context, code string) (*models.Invite, error)
	MarkUsed(ctx context.Context, id int64, usedBy int64) error
	DeletePendingByEmail(ctx context.Context, email string) error
}

type InviteRepository struct {
	db *pgxpool.Pool
}

func NewInviteRepository(db *pgxpool.Pool) *InviteRepository {
	return &InviteRepository{db: db}
}

func (r *InviteRepository) Create(ctx context.Context, invite *models.Invite) error {
	query := `INSERT INTO invites (code, email, created_by, expires_at)
	          VALUES ($1, $2, $3, $4)
	          RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		invite.Code,
		invite.Email,
		invite.CreatedBy,
		invite.ExpiresAt,
	).Scan(&invite.ID, &invite.CreatedAt)
}

func (r *InviteRepository) GetByCode(ctx context.Context, code string) (*models.Invite, error) {
	query := `SELECT id, code, email, created_by, used, used_by, expires_at, created_at
	          FROM invites WHERE code = $1`

	var inv models.Invite
	err := r.db.QueryRow(ctx, query, code).Scan(
		&inv.ID,
		&inv.Code,
		&inv.Email,
		&inv.CreatedBy,
		&inv.Used,
		&inv.UsedBy,
		&inv.ExpiresAt,
		&inv.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrInviteNotFound
	}
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *InviteRepository) MarkUsed(ctx context.Context, id int64, usedBy int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE invites SET used = TRUE, used_by = $1 WHERE id = $2`,
		usedBy, id,
	)
	return err
}

func (r *InviteRepository) DeletePendingByEmail(ctx context.Context, email string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM invites WHERE email = $1 AND used = FALSE`,
		email,
	)
	return err
}
