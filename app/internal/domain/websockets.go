package domain

import (
	"context"

	"golang.org/x/net/websocket"
)

type WebSocketHandler interface {
	Add(roomID int64, userID int64, conn *websocket.Conn)
	Remove(roomID int64, userID int64)
	LoggedIn(ctx context.Context, roomID int64, userID int64) error
	LoggedOut(ctx context.Context, roomID int64, userID int64) error
	GetOnlineUsers(ctx context.Context, roomID int64, userID int64) error
	HandleTrackEvent()
	PlayTrack(ctx context.Context, room *Room, track *Track, message WSRoomPlayTrackRequest) error
	PauseTrack(ctx context.Context, room *Room, user *User) error
	SeekTrack(ctx context.Context, room *Room, user *User, message WSRoomSeekTrackRequest) error
	SyncTrack(ctx context.Context, room *Room, user *User) error
	UpdateTrackTime(ctx context.Context, room *Room, user *User, message WSRoomUpdateTrackTimeRequest) error
}
