package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email,omitempty"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshToken struct {
	ID        int64     `json:"id"`
	SessionID string    `json:"session_id"`
	TokenHash string    `json:"-"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Invite struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Email     string    `json:"email"`
	CreatedBy int64     `json:"created_by"`
	Used      bool      `json:"used"`
	UsedBy    *int64    `json:"used_by,omitempty"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
