package repository

import "errors"

var (
	// Board game errors
	ErrBoardGameNotFound = errors.New("Board game not found")
	ErrDuplicateName     = errors.New("Board game with this name already exists")

	// Database errors
	ErrQueryFailed = errors.New("Database query failed")
)
