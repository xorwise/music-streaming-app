package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
	"github.com/xorwise/music-streaming-service/internal/utils"
)

type trackAddUsecase struct {
	trackRepository domain.TrackRepository
	roomRepository  domain.RoomRepository
	timeout         time.Duration
}

func NewTrackAddUsecase(tr domain.TrackRepository, rm domain.RoomRepository, timeout time.Duration) domain.TrackAddUsecase {
	return &trackAddUsecase{
		trackRepository: tr,
		roomRepository:  rm,
		timeout:         timeout,
	}
}

func (uc *trackAddUsecase) Create(ctx context.Context, track *domain.Track) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()
	return uc.trackRepository.Create(ctx, track)
}

func (uc *trackAddUsecase) FindAndSaveTrack(ctx context.Context, trackCh chan error, title string, artist string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()
	return utils.FindAndSaveTrack(ctx, trackCh, title, artist)
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

func (uc *trackAddUsecase) WaitForTrack(ctx context.Context, trackCh chan error, track *domain.Track) {
	err := <-trackCh
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), uc.timeout)
		defer cancel()
		uc.trackRepository.Remove(ctx, track)
		fmt.Println("deleted track")
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), uc.timeout)
		defer cancel()
		track.IsReady = true
		err := uc.trackRepository.Update(ctx, track)
		if err != nil {
			fmt.Println(err)
		}
	}
}
