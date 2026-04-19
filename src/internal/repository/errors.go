package repository

import "errors"

var (
	// Board game errors
	ErrBoardGameNotFound = errors.New("Board game not found")
	ErrDuplicateName     = errors.New("Board game with this name already exists")

	// Database errors
	ErrQueryFailed = errors.New("Database query failed")

	// Auth errors
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrSessionNotFound   = errors.New("session not found")
	ErrTokenNotFound     = errors.New("token not found")
	ErrInviteNotFound    = errors.New("invite not found")
)
