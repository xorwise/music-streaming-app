package usecase

import (
	"context"
	"errors"
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
		case domain.WSRoomPlayTrack:
			var req domain.WSRoomPlayTrackRequest
			req.TrackID = int64(message.Data.(map[string]interface{})["trackID"].(float64))
			req.Time = int64(message.Data.(map[string]interface{})["time"].(float64))
			track, err := wru.trackRepository.GetByID(context.Background(), req.TrackID)
			if err != nil {
				if errors.Is(err, domain.ErrTrackNotFound) {
					websocket.JSON.Send(ws, domain.WSRoomResponse{
						Type:  domain.WSRoomError,
						Data:  "",
						Error: "track not found",
					})
				}
				wru.log.Info(op, "error", err.Error(), "user", user.Username)
				break
			}

			if !track.IsReady {
				websocket.JSON.Send(ws, domain.WSRoomResponse{
					Type:  domain.WSRoomError,
					Data:  "",
					Error: "track not ready",
				})
				wru.log.Info(op, "error", "track not ready", "user", user.Username)
				break
			}

			err = wru.websocketHandler.PlayTrack(context.Background(), room, track, req)
			if err != nil {
				wru.log.Info(op, "error", err.Error(), "user", user.Username)
			}
		case domain.WSRoomPauseTrack:
			err := wru.websocketHandler.PauseTrack(context.Background(), room, user)
			if err != nil {
				wru.log.Info(op, "error", err.Error(), "user", user.Username)
			}
		case domain.WSRoomSeekTrack:
			var req domain.WSRoomSeekTrackRequest
			req.Time = int64(message.Data.(map[string]interface{})["time"].(float64))
			err := wru.websocketHandler.SeekTrack(context.Background(), room, user, req)
			if err != nil {
				wru.log.Info(op, "error", err.Error(), "user", user.Username)
			}
		case domain.WSRoomSyncTrack:
			err := wru.websocketHandler.SyncTrack(context.Background(), room, user)
			if err != nil {
				wru.log.Info(op, "error", err.Error(), "user", user.Username)
			}
		case domain.WSRoomUpdateTrackTime:
			var req domain.WSRoomUpdateTrackTimeRequest
			req.Time = int64(message.Data.(map[string]interface{})["time"].(float64))
			err := wru.websocketHandler.UpdateTrackTime(context.Background(), room, user, req)
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
