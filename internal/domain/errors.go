package domain

import "errors"

var (
	ErrNotFound       = errors.New("record not found")
	ErrEmailRequired  = errors.New("email is required")
	ErrInvalidEmail   = errors.New("invalid email address format")
	ErrDuplicateEmail = errors.New("email already subscribed")
	ErrDatabase       = errors.New("database operation failed")
)
