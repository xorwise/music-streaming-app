package domain

import "context"

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserCreateResponse struct {
	ID int64 `json:"id"`
}

type UserCreateUsecase interface {
	Create(ctx context.Context, user *User) (int64, error)
	HashPassword(password string) ([]byte, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}
