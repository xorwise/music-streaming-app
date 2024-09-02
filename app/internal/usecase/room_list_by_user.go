package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type roomListByUserUsecase struct {
	roomRepository domain.RoomRepository
	timeout        time.Duration
}

func NewRoomListByUserUsecase(rr domain.RoomRepository, timeout time.Duration) domain.RoomListByUserUsecase {
	return &roomListByUserUsecase{
		roomRepository: rr,
		timeout:        timeout,
	}
}

func (ru *roomListByUserUsecase) ListByUser(ctx context.Context, userID int64, limit int, offset int) ([]*domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.ListByUserID(ctx, userID, limit, offset)
}
