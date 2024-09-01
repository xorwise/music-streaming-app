package domain

import (
	"context"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
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
	CreateAccessToken(ctx context.Context, cfg *bootstrap.Config, user *User) (string, error)
}
