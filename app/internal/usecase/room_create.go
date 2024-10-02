package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type roomCreateUsecase struct {
	roomRepository domain.RoomRepository
	roomUtils      domain.RoomUtils
	timeout        time.Duration
}

func NewRoomCreateUsecase(rr domain.RoomRepository, ru domain.RoomUtils, timeout time.Duration) domain.RoomCreateUsecase {
	return &roomCreateUsecase{
		roomRepository: rr,
		roomUtils:      ru,
		timeout:        timeout,
	}
}

func (ru *roomCreateUsecase) Create(ctx context.Context, room *domain.Room) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.Create(ctx, room)
}

func (ru *roomCreateUsecase) GenerateCode(ctx context.Context, roomID int64) string {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomUtils.GenerateRoomCode(roomID)
}

func (ru *roomCreateUsecase) SetCode(ctx context.Context, roomID int64, code string) error {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.SetCode(ctx, roomID, code)
}

func (ru *roomCreateUsecase) AddRoomUser(ctx context.Context, roomID int64, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.AddRoomUser(ctx, roomID, userID)
}
