package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type roomUpdateAvatarUsecasee struct {
	roomRepository domain.RoomRepository
	roomUtils      domain.UserUtils
	timeout        time.Duration
}

func NewRoomUpdateAvatarUsecase(ur domain.RoomRepository, uu domain.UserUtils, timeout time.Duration) domain.RoomUpdateAvatarUsecase {
	return &roomUpdateAvatarUsecasee{
		roomRepository: ur,
		roomUtils:      uu,
		timeout:        timeout,
	}
}

func (ru *roomUpdateAvatarUsecasee) Update(ctx context.Context, room *domain.Room) error {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.Update(ctx, room)
}

func (ru *roomUpdateAvatarUsecasee) SaveFile(ctx context.Context, fileData string, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomUtils.SaveFile(ctx, fileData, filename)
}

func (ru *roomUpdateAvatarUsecasee) GetByID(ctx context.Context, id int64) (*domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.GetByID(ctx, id)
}

func (ru *roomUpdateAvatarUsecasee) GetByUserIDandRoomID(ctx context.Context, userID int64, roomID int64) (*domain.UserRoom, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	return ru.roomRepository.GetByUserIDandRoomID(ctx, userID, roomID)
}
