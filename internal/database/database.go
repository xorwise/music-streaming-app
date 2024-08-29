package database

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrUserExists = errors.New("user already exists")
)
