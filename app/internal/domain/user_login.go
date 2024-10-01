package domain

import (
	"context"
)

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	AccessToken string `json:"access_token"`
}

type UserLoginUsecase interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	CreateAccessToken(ctx context.Context, user *User) (string, error)
	CheckPasswordHash(password string, hash string) bool
}
