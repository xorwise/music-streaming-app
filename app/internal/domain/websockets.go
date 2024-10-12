package domain

import (
	"context"
	"database/sql"

	"golang.org/x/net/websocket"
)

type WSClients map[int64]map[int64]*websocket.Conn

type WebSocketHandler interface {
	Add(roomID int64, userID int64, conn *websocket.Conn)
	Remove(roomID int64, userID int64)
	LoggedIn(ctx context.Context, roomID int64, userID int64) (*WSRoomResponse, error)
	LoggedOut(ctx context.Context, roomID int64, userID int64) (*WSRoomResponse, error)
	GetOnlineUsers(ctx context.Context, roomID int64, userID int64, additionalClients []int64) error
	HandleTrackEvent()
	PlayTrack(ctx context.Context, room *Room, track *Track, message WSRoomPlayTrackRequest) (*WSRoomResponse, error)
	PauseTrack(ctx context.Context, room *Room, user *User) (*WSRoomResponse, error)
	SeekTrack(ctx context.Context, room *Room, user *User, message WSRoomSeekTrackRequest) (*WSRoomResponse, error)
	SyncTrack(ctx context.Context, room *Room, user *User) (*WSRoomResponse, error)
	UpdateTrackTime(ctx context.Context, room *Room, user *User, message WSRoomUpdateTrackTimeRequest) error
	StopTrack(ctx context.Context, room *Room, user *User) (*WSRoomResponse, error)
	BroadcastClients(roomID int64) []int64
	BroadcastMsg(broadcast chan *RoomBroadcastResponse, db *sql.DB)
}

type RoomBroadcastResponse struct {
	RoomID   int64
	Response *WSRoomResponse
}
