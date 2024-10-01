package domain

import (
	"context"

	"golang.org/x/net/websocket"
)

const (
	WSRoomError = iota
	WSRoomLoggedIn
	WSRoomLoggedOut
	WSRoomGetOnlineUsers
	WSRoomFetchMusicChunks
	WSRoomTrackEvent
	WSRoomPlayTrack
	WSRoomPauseTrack
	WSRoomSeekTrack
	WSRoomSyncTrack
	WSRoomUpdateTrackTime
	WSRoomLoggedInTrack
	WSRoomStopTrack
)

type WSRoomUsecase interface {
	Handle(ws *websocket.Conn, room *Room, user *User)
	GetByID(ctx context.Context, id int64) (*Room, error)
	GetUserIDandRoomID(ctx context.Context, roomID int64, userID int64) (*UserRoom, error)
	LoggedIn(ctx context.Context, roomID int64, userID int64, ws *websocket.Conn)
	LoggedOut(ctx context.Context, roomID int64, userID int64)
}

type WSRoomRequest struct {
	Type int         `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

type WSRoomResponse struct {
	Type  int         `json:"type"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

type WSRoomPlayTrackRequest struct {
	TrackID int64 `json:"trackID"`
	Time    int64 `json:"time"`
}

type WSRoomSeekTrackRequest struct {
	Time int64 `json:"time"`
}

type WSRoomUpdateTrackTimeRequest struct {
	Time int64 `json:"time"`
}
