package usecase

import (
	"context"
	"net/url"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type trackListByRoomUsecase struct {
	trackRepository domain.TrackRepository
	roomRepository  domain.RoomRepository
	timeout         time.Duration
}

func NewTrackListByRoomUsecase(tr domain.TrackRepository, rm domain.RoomRepository, timeout time.Duration) domain.TrackListByRoomUsecase {
	return &trackListByRoomUsecase{
		trackRepository: tr,
		roomRepository:  rm,
		timeout:         timeout,
	}
}

func (tu *trackListByRoomUsecase) GetRoomByID(ctx context.Context, id int64) (*domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.timeout)
	defer cancel()
	return tu.roomRepository.GetByID(ctx, id)
}

func (tu *trackListByRoomUsecase) GetByUserIDandRoomID(ctx context.Context, roomID int64, userID int64) (*domain.UserRoom, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.timeout)
	defer cancel()
	return tu.roomRepository.GetByUserIDandRoomID(ctx, roomID, userID)
}

func (tu *trackListByRoomUsecase) ListByRoomID(ctx context.Context, roomID int64, params url.Values) ([]*domain.Track, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.timeout)
	defer cancel()
	return tu.trackRepository.ListByRoomID(ctx, roomID, params)
}
