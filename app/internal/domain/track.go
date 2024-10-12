package domain

import (
	"context"
	"net/url"
)

type Track struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	Path    string `json:"path"`
	RoomID  int64  `json:"room_id"`
	IsReady bool   `json:"is_ready"`
}

type TrackRepository interface {
	Create(ctx context.Context, track *Track) (int64, error)
	GetByID(ctx context.Context, trackID int64) (*Track, error)
	Remove(ctx context.Context, track *Track) error
	Update(ctx context.Context, track *Track) error
	ListByRoomID(ctx context.Context, roomID int64, params url.Values) ([]*Track, error)
	RemoveOutdated(tu TrackUtils)
}

type TrackUtils interface {
	FindAndSaveTrack(ctx context.Context, trackCh chan error, title string, artist string) (string, error)
	RemoveFiles(ctx context.Context, track *Track) error
}

type TrackStatus struct {
	ID      int64
	RoomID  int64
	Path    string
	IsReady bool
}
