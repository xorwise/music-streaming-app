package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type roomLeaveUsecase struct {
	roomRepository domain.RoomRepository
	timeout        time.Duration
}

func NewRoomLeaveUsecase(rr domain.RoomRepository, timeout time.Duration) domain.RoomLeaveUsecase {
	return &roomLeaveUsecase{
		roomRepository: rr,
		timeout:        timeout,
	}
}

func (ru *roomLeaveUsecase) GetByID(ctx context.Context, id int64) (*domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.GetByID(ctx, id)
}

func (ru *roomLeaveUsecase) RemoveRoomUser(ctx context.Context, roomID int64, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.RemoveRoomUser(ctx, roomID, userID)
}

func (ru *roomLeaveUsecase) GetUserIDandRoomID(ctx context.Context, id int64, userID int64) (*domain.UserRoom, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.GetByUserIDandRoomID(ctx, id, userID)
}
