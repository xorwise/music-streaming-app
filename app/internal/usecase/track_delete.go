package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type trackDeleteUsecase struct {
	trackRepository domain.TrackRepository
	roomRepository  domain.RoomRepository
	trackUtils      domain.TrackUtils
	timeout         time.Duration
}

func NewTrackDeleteUsecase(
	tr domain.TrackRepository,
	rr domain.RoomRepository,
	tu domain.TrackUtils,
	timeout time.Duration,
) domain.TrackDeleteUsecase {
	return &trackDeleteUsecase{
		trackRepository: tr,
		roomRepository:  rr,
		trackUtils:      tu,
		timeout:         timeout,
	}
}

func (tu *trackDeleteUsecase) GetByID(ctx context.Context, id int64) (*domain.Track, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.timeout)
	defer cancel()
	return tu.trackRepository.GetByID(ctx, id)
}

func (tu *trackDeleteUsecase) GetByUserIDandRoomID(ctx context.Context, roomID int64, userID int64) (*domain.UserRoom, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.timeout)
	defer cancel()
	return tu.roomRepository.GetByUserIDandRoomID(ctx, roomID, userID)
}

func (tu *trackDeleteUsecase) Remove(ctx context.Context, track *domain.Track) error {
	ctx, cancel := context.WithTimeout(ctx, tu.timeout)
	defer cancel()
	return tu.trackRepository.Remove(ctx, track)
}

func (tu *trackDeleteUsecase) RemoveFiles(ctx context.Context, track *domain.Track) error {
	ctx, cancel := context.WithTimeout(ctx, tu.timeout)
	defer cancel()
	return tu.trackUtils.RemoveFiles(ctx, track)
}
