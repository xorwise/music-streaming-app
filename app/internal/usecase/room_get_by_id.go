package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type roomGetByIDUsecase struct {
	roomRepository domain.RoomRepository
	timeout        time.Duration
}

func NewRoomGetByIDUsecase(rr domain.RoomRepository, timeout time.Duration) domain.RoomGetByIDUsecase {
	return &roomGetByIDUsecase{
		roomRepository: rr,
		timeout:        timeout,
	}
}

func (ru *roomGetByIDUsecase) GetByID(ctx context.Context, id int64) (*domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.GetByID(ctx, id)
}

func (ru *roomGetByIDUsecase) GetUserIDandRoomID(ctx context.Context, id int64, userID int64) (*domain.UserRoom, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.GetByUserIDandRoomID(ctx, id, userID)
}
