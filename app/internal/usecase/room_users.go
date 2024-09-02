package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type roomUsersUsecase struct {
	roomRepository domain.RoomRepository
	timeout        time.Duration
}

func NewRoomUsersUsecase(rr domain.RoomRepository, timeout time.Duration) domain.RoomUsersUsecase {
	return &roomUsersUsecase{
		roomRepository: rr,
		timeout:        timeout,
	}
}

func (ru *roomUsersUsecase) ListRoomUsers(ctx context.Context, roomID int64, limit int, offset int) ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.ListRoomUsers(ctx, roomID, limit, offset)
}

func (ru *roomUsersUsecase) GetByUserIDandRoomID(ctx context.Context, id int64, userID int64) (*domain.UserRoom, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.GetByUserIDandRoomID(ctx, id, userID)
}

func (ru *roomUsersUsecase) GetByID(ctx context.Context, id int64) (*domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.GetByID(ctx, id)
}
