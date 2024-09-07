package usecase

import (
	"context"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
	"golang.org/x/net/websocket"
)

type wsRoomUsecase struct {
	roomRepository domain.RoomRepository
	timeout        time.Duration
	ws             *websocket.Conn
	clients        *domain.WSClients
}

func NewWSRoomUsecase(rr domain.RoomRepository, timeout time.Duration, clients *domain.WSClients) domain.WSRoomUsecase {
	return &wsRoomUsecase{
		roomRepository: rr,
		timeout:        timeout,
		clients:        clients,
	}
}

func (wru *wsRoomUsecase) Handle(ws *websocket.Conn, room *domain.Room, user *domain.User) {
	var message domain.WSRoomRequest

	for {
		err := websocket.JSON.Receive(ws, &message)
		if err != nil {
			websocket.JSON.Send(ws, domain.WSRoomResponse{
				Type:  0,
				Data:  "",
				Error: err.Error(),
			})
			break
		}

		switch message.Type {
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
	}

	for _, client := range clients {
		websocket.JSON.Send(client, response)
	}
}

func (wru *wsRoomUsecase) LoggedOut(ctx context.Context, roomID int64, userID int64) {
	ctx, cancel := context.WithTimeout(ctx, wru.timeout)
	defer cancel()

	var response domain.WSRoomResponse
	response.Type = 1
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
	}

	for _, client := range clients {
		websocket.JSON.Send(client, response)
	}
}
