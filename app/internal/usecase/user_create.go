package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type userCreateUsecase struct {
	userRepository domain.UserRepository
	userUtils      domain.UserUtils
	timeout        time.Duration
}

func NewUserCreateUsecase(ur domain.UserRepository, uu domain.UserUtils, timeout time.Duration) domain.UserCreateUsecase {
	return &userCreateUsecase{
		userRepository: ur,
		userUtils:      uu,
		timeout:        timeout,
	}
}

func (uu *userCreateUsecase) Create(ctx context.Context, user *domain.User) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()
	return uu.userRepository.Create(ctx, user)
}

func (uu *userCreateUsecase) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()
	return uu.userRepository.GetByUsername(ctx, username)
}

func (uu *userCreateUsecase) HashPassword(password string) ([]byte, error) {
	return uu.userUtils.HashPassword(password)
}
