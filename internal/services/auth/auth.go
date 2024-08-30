package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain/auth/models"
	"github.com/xorwise/music-streaming-service/internal/lib/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	SaveUser(ctx context.Context, username string, passHash []byte) (userID int64, err error)
	User(ctx context.Context, username string) (models.User, error)
}

type Auth struct {
	log      *slog.Logger
	us       UserStorage
	tokenTTL time.Duration
}

func New(log *slog.Logger, us UserStorage, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:      log,
		us:       us,
		tokenTTL: tokenTTL,
	}
}

func (a *Auth) Register(ctx context.Context, username string, password string) (userID int64, err error) {
	const op = "Auth.Register"

	log := a.log.With(slog.String("op", op), slog.String("username", username))

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.us.SaveUser(ctx, username, passHash)
	if err != nil {
		log.Error("failed to save user", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) Login(ctx context.Context, username string, password string) (token string, err error) {
	const op = "Auth.Login"

	log := a.log.With(slog.String("op", op), slog.String("username", username))

	log.Info("logging in user")

	user, err := a.us.User(ctx, username)
	if err != nil {
		log.Error("failed to get user", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		log.Error("failed to compare passwords", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err = jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		log.Error("failed to create token", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
