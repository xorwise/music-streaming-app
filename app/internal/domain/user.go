package domain

import "context"

type User struct {
	ID       int64
	Username string
	Avatar   string
	PassHash string
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (int64, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type UserUtils interface {
	HashPassword(password string) ([]byte, error)
	CheckPasswordHash(password, hash string) bool
	CreateAccessToken(ctx context.Context, user *User) (string, error)
	SaveFile(ctx context.Context, fileData string, filename string) (string, error)
}
