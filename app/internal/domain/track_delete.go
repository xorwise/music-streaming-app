package domain

import "context"

type TrackDeleteUsecase interface {
	GetByID(ctx context.Context, id int64) (*Track, error)
	GetByUserIDandRoomID(ctx context.Context, roomID int64, userID int64) (*UserRoom, error)
	Remove(ctx context.Context, track *Track) error
	RemoveFiles(ctx context.Context, track *Track) error
}
