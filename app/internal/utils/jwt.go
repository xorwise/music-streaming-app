package utils

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

func CreateAccessToken(ctx context.Context, cfg *bootstrap.Config, user *domain.User) (string, error) {
	exp := time.Now().Add(time.Duration(cfg.TokenTTL) * time.Second)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = exp.Unix()

	tokenStr, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
