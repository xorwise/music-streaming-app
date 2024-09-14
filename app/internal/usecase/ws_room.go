package usecase

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
	"github.com/xorwise/music-streaming-service/internal/utils"
	"golang.org/x/net/websocket"
)

// TODO: add other layer for websocket handling

type wsRoomUsecase struct {
	roomRepository  domain.RoomRepository
	trackRepository domain.TrackRepository
	timeout         time.Duration
	ws              *websocket.Conn
	clients         *domain.WSClients
	log             *slog.Logger
	trackCh         chan domain.TrackStatus
}

func NewWSRoomUsecase(
	rr domain.RoomRepository,
	tr domain.TrackRepository,
	timeout time.Duration,
	clients *domain.WSClients,
	log *slog.Logger,
	trackCh chan domain.TrackStatus,
) domain.WSRoomUsecase {
	return &wsRoomUsecase{
		roomRepository:  rr,
		trackRepository: tr,
		timeout:         timeout,
		clients:         clients,
		log:             log,
		trackCh:         trackCh,
	}
}

func (wru *wsRoomUsecase) Handle(ws *websocket.Conn, room *domain.Room, user *domain.User) {
	var message domain.WSRoomRequest

	const op = "Websockets.Room"

	go wru.handleTrackEvent()

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
		case 4:
			trackID, ok := message.Data.(float64)
			if !ok {
				websocket.JSON.Send(ws, domain.WSRoomResponse{
					Type:  0,
					Data:  "",
					Error: "data is not int",
				})
				wru.log.Info(op, "error", "data is not int64", "user", user.Username)
				break
			}
			wru.FetchMusicChunks(context.Background(), int64(trackID), room.ID, user.ID)
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

func (wru *wsRoomUsecase) FetchMusicChunks(ctx context.Context, trackID int64, roomID int64, userID int64) {
	ctx, cancel := context.WithTimeout(ctx, wru.timeout)
	defer cancel()

	track, err := wru.trackRepository.GetByID(ctx, trackID)
	if err != nil {
		websocket.JSON.Send(wru.clients.RoomClients[roomID][userID], domain.WSRoomResponse{
			Type:  0,
			Data:  "",
			Error: "track not found",
		})
		wru.log.Error("Websockets.Room.FetchMusicChunks", "error", err.Error(), "user", userID)
		return
	}

	const op = "Websockets.Room.FetchMusicChunks"

	musicReader, err := utils.NewMusicReader(track.Path)
	if err != nil {
		websocket.JSON.Send(wru.clients.RoomClients[roomID][userID], domain.WSRoomResponse{
			Type:  0,
			Data:  "",
			Error: err.Error(),
		})
		wru.log.Error(op, "error", err.Error(), "user", userID)
		return
	}
	defer musicReader.Close()

	buff := make([]byte, 6*1024)

	for {
		n, err := musicReader.Read(buff)
		if err == io.EOF {
			return
		}
		if err != nil {
			wru.log.Error(op, "error", err.Error(), "user", userID)
			return
		}

		websocket.Message.Send(wru.clients.RoomClients[roomID][userID], buff[:n])
	}
}

func (wru *wsRoomUsecase) handleTrackEvent() {
	fmt.Println("working...")
	trackEvent := <-wru.trackCh
	var event string
	if trackEvent.IsReady {
		event = "ready"

	} else {
		event = "removed"
	}
	msg := domain.WSRoomResponse{
		Type: 5,
		Data: struct {
			TrackID int64  `json:"trackID"`
			Event   string `json:"event"`
		}{
			TrackID: trackEvent.ID,
			Event:   event,
		},
	}
	for _, client := range wru.clients.RoomClients[trackEvent.RoomID] {
		websocket.JSON.Send(client, msg)
	}
}
