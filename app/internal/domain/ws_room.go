package domain

import (
	"context"

	"golang.org/x/net/websocket"
)

type WSRoomUsecase interface {
	Handle(ws *websocket.Conn, room *Room, user *User)
	LoggedIn(ctx context.Context, roomID int64, userID int64)
	GetByID(ctx context.Context, id int64) (*Room, error)
	GetUserIDandRoomID(ctx context.Context, id int64, userID int64) (*UserRoom, error)
	LoggedOut(ctx context.Context, roomID int64, userID int64)
	GetOnlineUsers(ctx context.Context, roomID int64, userID int64)
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
