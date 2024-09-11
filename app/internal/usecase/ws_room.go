package usecase

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
	"golang.org/x/net/websocket"
)

type wsRoomUsecase struct {
	roomRepository domain.RoomRepository
	timeout        time.Duration
	ws             *websocket.Conn
	clients        *domain.WSClients
	log            *slog.Logger
}

func NewWSRoomUsecase(rr domain.RoomRepository, timeout time.Duration, clients *domain.WSClients, log *slog.Logger) domain.WSRoomUsecase {
	return &wsRoomUsecase{
		roomRepository: rr,
		timeout:        timeout,
		clients:        clients,
		log:            log,
	}
}

func (wru *wsRoomUsecase) Handle(ws *websocket.Conn, room *domain.Room, user *domain.User) {
	var message domain.WSRoomRequest

	const op = "Websockets.Room"

	for {
		err := websocket.JSON.Receive(ws, &message)
		if err != nil {
			websocket.JSON.Send(ws, domain.WSRoomResponse{
				Type:  0,
				Data:  "",
				Error: err.Error(),
			})
			wru.log.Info(op, "error", err.Error(), "user", user.Username)
			break
		}

		wru.log.Info(op, "received messae with type", message.Type, "user", user.Username)

		switch message.Type {
		case 3:
			wru.GetOnlineUsers(context.Background(), room.ID, user.ID)
		}
	}

}

func (wru *wsRoomUsecase) GetByID(ctx context.Context, id int64) (*domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, wru.timeout)
	defer cancel()
	return wru.roomRepository.GetByID(ctx, id)
}

func (wru *wsRoomUsecase) GetUserIDandRoomID(ctx context.Context, id int64, userID int64) (*domain.UserRoom, error) {
	ctx, cancel := context.WithTimeout(ctx, wru.timeout)
	defer cancel()
	return wru.roomRepository.GetByUserIDandRoomID(ctx, id, userID)
}

func (wru *wsRoomUsecase) LoggedIn(ctx context.Context, roomID int64, userID int64) {
	ctx, cancel := context.WithTimeout(ctx, wru.timeout)
	defer cancel()

	const op = "Websockets.Room.LoggedIn"

	var response domain.WSRoomResponse
	response.Type = 1
	response.Data = struct {
		Event  string `json:"event"`
		UserID int64  `json:"userID"`
	}{
		Event:  "connected",
		UserID: userID,
	}

	clients, ok := wru.clients.RoomClients[roomID]
	if !ok {
		websocket.JSON.Send(wru.ws, domain.WSRoomResponse{
			Type:  0,
			Data:  "",
			Error: "internal server error",
		})
		wru.log.Error(op, "error", "internl server error", "user", userID)
		return
	}

	for _, client := range clients {
		websocket.JSON.Send(client, response)
	}
}

func (wru *wsRoomUsecase) LoggedOut(ctx context.Context, roomID int64, userID int64) {
	ctx, cancel := context.WithTimeout(ctx, wru.timeout)
	defer cancel()

	const op = "Websockets.Room.LoggedOut"

	var response domain.WSRoomResponse
	response.Type = 2
	response.Data = struct {
		Event  string `json:"event"`
		UserID int64  `json:"userID"`
	}{
		Event:  "disconnected",
		UserID: userID,
	}

	clients, ok := wru.clients.RoomClients[roomID]
	if !ok {
		websocket.JSON.Send(wru.ws, domain.WSRoomResponse{
			Type:  0,
			Data:  "",
			Error: "internal server error",
		})
		wru.log.Error(op, "error", "internl server error", "user", userID)
		return
	}

	for _, client := range clients {
		websocket.JSON.Send(client, response)
	}
}

func (wru *wsRoomUsecase) GetOnlineUsers(ctx context.Context, roomID int64, userID int64) {
	ctx, cancel := context.WithTimeout(ctx, wru.timeout)
	defer cancel()

	const op = "Websockets.Room.GetOnlineUsers"

	clients, ok := wru.clients.RoomClients[roomID]
	if !ok {
		websocket.JSON.Send(wru.ws, domain.WSRoomResponse{
			Type:  0,
			Data:  "",
			Error: "internal server error",
		})
		wru.log.Error(op, "error", "internl server error", "user", userID)
		return
	}

	var response domain.WSRoomResponse
	response.Type = 3
	response.Data = make([]int64, 0)

	for id := range clients {
		response.Data = append(response.Data.([]int64), id)
	}

	for id, client := range clients {
		if id == userID {
			websocket.JSON.Send(client, response)
			break
		}
	}
}