package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type roomEnterUsecase struct {
	roomRepository domain.RoomRepository
	timeout        time.Duration
}

func NewRoomEnterUsecase(rr domain.RoomRepository, timeout time.Duration) domain.RoomEnterUsecase {
	return &roomEnterUsecase{
		roomRepository: rr,
		timeout:        timeout,
	}
}

func (ru *roomEnterUsecase) GetByCode(ctx context.Context, code string) (*domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.GetByCode(ctx, code)
}

func (ru *roomEnterUsecase) GetByUserIDandRoomID(ctx context.Context, id int64, userID int64) (*domain.UserRoom, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.GetByUserIDandRoomID(ctx, id, userID)
}

func (ru *roomEnterUsecase) AddRoomUser(ctx context.Context, roomID int64, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.AddRoomUser(ctx, roomID, userID)
}
