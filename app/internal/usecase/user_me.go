package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type userMeUsecase struct {
	userRepository domain.UserRepository
	timeout        time.Duration
}

func NewUserMeUsecase(ur domain.UserRepository, timeout time.Duration) domain.UserMeUsecase {
	return &userMeUsecase{
		userRepository: ur,
		timeout:        timeout,
	}
}

func (uu *userMeUsecase) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()
	return uu.userRepository.GetByID(ctx, id)
}
