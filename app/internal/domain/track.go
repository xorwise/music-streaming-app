package domain

import (
	"context"
	"net/url"
)

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
	Remove(ctx context.Context, track *Track) error
	Update(ctx context.Context, track *Track) error
	ListByRoomID(ctx context.Context, roomID int64, params url.Values) ([]*Track, error)
}

type TrackStatus struct {
	ID      int64
	RoomID  int64
	IsReady bool
}
