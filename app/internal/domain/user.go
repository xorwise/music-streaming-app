package domain

import "context"

type User struct {
	ID       int64
	Username string
	PassHash string
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (int64, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}
