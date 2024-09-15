package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
	"golang.org/x/net/websocket"
)

type wsRoomUsecase struct {
	roomRepository   domain.RoomRepository
	trackRepository  domain.TrackRepository
	websocketHandler domain.WebSocketHandler
	timeout          time.Duration
	log              *slog.Logger
}

func NewWSRoomUsecase(
	rr domain.RoomRepository,
	tr domain.TrackRepository,
	wsh domain.WebSocketHandler,
	log *slog.Logger,
	timeout time.Duration,
) domain.WSRoomUsecase {
	return &wsRoomUsecase{
		roomRepository:   rr,
		trackRepository:  tr,
		websocketHandler: wsh,
		timeout:          timeout,
		log:              log,
	}
}

func (wru *wsRoomUsecase) Handle(ws *websocket.Conn, room *domain.Room, user *domain.User) {
	var message domain.WSRoomRequest

	const op = "Websockets.Room"

	for {
		err := websocket.JSON.Receive(ws, &message)
		if err != nil {
			err1 := websocket.JSON.Send(ws, domain.WSRoomResponse{
				Type:  domain.WSRoomError,
				Data:  "",
				Error: err.Error(),
			})
			if err1 != nil {
				break
			}
			wru.log.Info(op, "error", err.Error(), "user", user.Username)
			break
		}

		wru.log.Info(op, "received messae with type", message.Type, "user", user.Username)

		switch message.Type {
		case domain.WSRoomGetOnlineUsers:
			err := wru.websocketHandler.GetOnlineUsers(context.Background(), room.ID, user.ID)
			if err != nil {
				wru.log.Info(op, "error", err.Error(), "user", user.Username)
			}
		case domain.WSRoomFetchMusicChunks:
			trackID, ok := message.Data.(float64)
			if !ok {
				websocket.JSON.Send(ws, domain.WSRoomResponse{
					Type:  domain.WSRoomError,
					Data:  "",
					Error: "data is not int",
				})
				wru.log.Info(op, "error", "data is not int64", "user", user.Username)
				break
			}
			ctx, cancel := context.WithTimeout(context.Background(), wru.timeout)
			defer cancel()
			track, err := wru.trackRepository.GetByID(ctx, int64(trackID))
			if err != nil {
				websocket.JSON.Send(ws, domain.WSRoomResponse{
					Type:  domain.WSRoomError,
					Data:  "",
					Error: err.Error(),
				})
				wru.log.Info(op, "error", err.Error(), "user", user.Username)
				break
			}
			err = wru.websocketHandler.FetchMusicChunks(context.Background(), track, room.ID, user.ID)
			if err != nil {
				wru.log.Info(op, "error", err.Error(), "user", user.Username)
			}
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

func (wru *wsRoomUsecase) LoggedIn(ctx context.Context, roomID int64, userID int64, ws *websocket.Conn) {
	ctx, cancel := context.WithTimeout(ctx, wru.timeout)
	defer cancel()

	wru.websocketHandler.Add(roomID, userID, ws)
	err := wru.websocketHandler.LoggedIn(ctx, roomID, userID)
	if err != nil {
		wru.log.Info("Websockets.Room", "error", err.Error(), "user", userID)
	}
}

func (wru *wsRoomUsecase) LoggedOut(ctx context.Context, roomID int64, userID int64) {
	ctx, cancel := context.WithTimeout(ctx, wru.timeout)
	defer cancel()
	wru.websocketHandler.Remove(roomID, userID)
	err := wru.websocketHandler.LoggedOut(ctx, roomID, userID)
	if err != nil {
		wru.log.Info("Websockets.Room", "error", err.Error(), "user", userID)
	}
}
