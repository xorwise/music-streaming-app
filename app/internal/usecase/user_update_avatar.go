package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type userUpdateAvatarUsecase struct {
	userRepository domain.UserRepository
	userUtils      domain.UserUtils
	timeout        time.Duration
}

func NewUserUpdateAvatarUsecase(ur domain.UserRepository, uu domain.UserUtils, timeout time.Duration) domain.UserUpdateAvatarUsecase {
	return &userUpdateAvatarUsecase{
		userRepository: ur,
		userUtils:      uu,
		timeout:        timeout,
	}
}

func (uu *userUpdateAvatarUsecase) Update(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()
	return uu.userRepository.Update(ctx, user)
}

func (uu *userUpdateAvatarUsecase) SaveFile(ctx context.Context, fileData string, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()
	return uu.userUtils.SaveFile(ctx, fileData, filename)
}
