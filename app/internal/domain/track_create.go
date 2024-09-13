package domain

import "context"

type TrackAddRequest struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	RoomID int64  `json:"room_id"`
}

type TrackAddResponse struct {
	ID int64 `json:"id"`
}

type TrackAddUsecase interface {
	Create(ctx context.Context, track *Track) (int64, error)
	FindAndSaveTrack(ctx context.Context, trackCh chan error, title string, artist string) (string, error)
	GetRoomByID(ctx context.Context, id int64) (*Room, error)
	GetByUserIDandRoomID(ctx context.Context, roomID int64, userID int64) (*UserRoom, error)
	WaitForTrack(ctx context.Context, trackCh chan error, track *Track)
}
