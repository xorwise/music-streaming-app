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
	FetchMusicChunks(ctx context.Context, track *Track, roomID int64, userID int64) error
	HandleTrackEvent()
}
