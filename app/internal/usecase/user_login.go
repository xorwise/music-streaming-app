package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type userLoginUsecase struct {
	userRepository domain.UserRepository
	userUtils      domain.UserUtils
	timeout        time.Duration
}

func NewUserLoginUsecase(ur domain.UserRepository, uu domain.UserUtils, timeout time.Duration) domain.UserLoginUsecase {
	return &userLoginUsecase{
		userRepository: ur,
		userUtils:      uu,
		timeout:        timeout,
	}
}

func (uu *userLoginUsecase) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()
	return uu.userRepository.GetByUsername(ctx, username)
}

func (uu *userLoginUsecase) CreateAccessToken(ctx context.Context, user *domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()
	return uu.userUtils.CreateAccessToken(ctx, user)
}

func (uu *userLoginUsecase) CheckPasswordHash(password string, hash string) bool {
	return uu.userUtils.CheckPasswordHash(password, hash)
}
