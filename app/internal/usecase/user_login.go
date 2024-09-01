package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
	"github.com/xorwise/music-streaming-service/internal/utils"
)

type userLoginUsecase struct {
	userRepository domain.UserRepository
	timeout        time.Duration
}

func NewUserLoginUsecase(ur domain.UserRepository, timeout time.Duration) domain.UserLoginUsecase {
	return &userLoginUsecase{
		userRepository: ur,
		timeout:        timeout,
	}
}

func (uu *userLoginUsecase) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()
	return uu.userRepository.GetByUsername(ctx, username)
}

func (uu *userLoginUsecase) CreateAccessToken(ctx context.Context, cfg *bootstrap.Config, user *domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()
	return utils.CreateAccessToken(ctx, cfg, user)
}
