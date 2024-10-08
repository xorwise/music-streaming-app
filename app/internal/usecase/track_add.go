package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type trackAddUsecase struct {
	trackRepository domain.TrackRepository
	roomRepository  domain.RoomRepository
	trackUtils      domain.TrackUtils
	chanErr         chan error
	timeout         time.Duration
}

func NewTrackAddUsecase(tr domain.TrackRepository, rm domain.RoomRepository, tu domain.TrackUtils, chanErr chan error, timeout time.Duration) domain.TrackAddUsecase {
	return &trackAddUsecase{
		trackRepository: tr,
		roomRepository:  rm,
		trackUtils:      tu,
		chanErr:         chanErr,
		timeout:         timeout,
	}
}

func (uc *trackAddUsecase) Create(ctx context.Context, track *domain.Track) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()
	return uc.trackRepository.Create(ctx, track)
}

func (uc *trackAddUsecase) FindAndSaveTrack(ctx context.Context, title string, artist string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()
	return uc.trackUtils.FindAndSaveTrack(ctx, uc.chanErr, title, artist)
}

func (uc *trackAddUsecase) GetRoomByID(ctx context.Context, id int64) (*domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()
	return uc.roomRepository.GetByID(ctx, id)
}

func (uc *trackAddUsecase) GetByUserIDandRoomID(ctx context.Context, roomID int64, userID int64) (*domain.UserRoom, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()
	return uc.roomRepository.GetByUserIDandRoomID(ctx, roomID, userID)
}

func (uc *trackAddUsecase) WaitForTrack(ctx context.Context, track *domain.Track) {
	err := <-uc.chanErr
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), uc.timeout)
		defer cancel()
		uc.trackRepository.Remove(ctx, track)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), uc.timeout)
		defer cancel()
		track.IsReady = true
		uc.trackRepository.Update(ctx, track)
	}
}
