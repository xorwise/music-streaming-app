package domain

import "context"

type Track struct {
	ID      int64
	Title   string
	Artist  string
	Path    string
	RoomID  int64
	IsReady bool
}

type TrackRepository interface {
	Create(ctx context.Context, track *Track) (int64, error)
	GetByID(ctx context.Context, trackID int64) (*Track, error)
	Remove(ctx context.Context, trackID int64) error
	Update(ctx context.Context, track *Track) error
}
