package domain

import (
	"context"
	"net/url"
)

type TrackListByRoomUsecase interface {
	GetRoomByID(ctx context.Context, id int64) (*Room, error)
	GetByUserIDandRoomID(ctx context.Context, roomID int64, userID int64) (*UserRoom, error)
	ListByRoomID(ctx context.Context, roomID int64, params url.Values) ([]*Track, error)
}

type TrackListByRoomResponse struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	RoomID  int64  `json:"room_id"`
	Path    string `json:"path"`
	IsReady bool   `json:"is_ready"`
}
