package usecase

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
	"golang.org/x/net/websocket"
)

type wsRoomUsecase struct {
	roomRepository   domain.RoomRepository
	trackRepository  domain.TrackRepository
	websocketHandler domain.WebSocketHandler
	msgBrokerUtil    domain.MessageBrokerUtils
	log              *slog.Logger
	prom             *bootstrap.Prometheus
	timeout          time.Duration
}

func NewWSRoomUsecase(
	rr domain.RoomRepository,
	tr domain.TrackRepository,
	wsh domain.WebSocketHandler,
	mbu domain.MessageBrokerUtils,
	log *slog.Logger,
	prom *bootstrap.Prometheus,
	timeout time.Duration,
) domain.WSRoomUsecase {
	return &wsRoomUsecase{
		roomRepository:   rr,
		trackRepository:  tr,
		websocketHandler: wsh,
		msgBrokerUtil:    mbu,
		log:              log,
		prom:             prom,
		timeout:          timeout,
	}
}

func (wru *wsRoomUsecase) Handle(ws *websocket.Conn, room *domain.Room, user *domain.User) {
	var message domain.WSRoomRequest

	const op = "Websockets.Room"

	for {
		err := websocket.JSON.Receive(ws, &message)
		startTime := time.Now()
		if err != nil {
			err1 := websocket.JSON.Send(ws, domain.WSRoomResponse{
				Type:  domain.WSRoomError,
				Data:  "",
				Error: err.Error(),
			})
			duration := time.Since(startTime).Seconds()
			wru.prom.WebsocketMessageHandlingDuration.WithLabelValues("ws/room?id="+strconv.Itoa(int(room.ID)), "0", strconv.Itoa(int(user.ID))).Observe(duration)
			if err1 != nil {
				break
			}
			wru.log.Info(op, "error", err.Error(), "user", user.Username)
			break
		}

		wru.log.Info(op, "received messae with type", message.Type, "user", user.Username)

		switch message.Type {
		case domain.WSRoomGetOnlineUsers:
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			clients, err := wru.msgBrokerUtil.GetClientsInRoom(ctx, room.ID)
			if err != nil {
				wru.log.Info(op, "error", err.Error(), "user", user.Username)
				break
			}
			err = wru.websocketHandler.GetOnlineUsers(context.Background(), room.ID, user.ID, clients)
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

			response, err := wru.websocketHandler.PlayTrack(context.Background(), room, track, req)
			wru.msgBrokerUtil.BroadcastMessage(room.ID, response)
			if err != nil {
				wru.log.Info(op, "error", err.Error(), "user", user.Username)
			}
		case domain.WSRoomPauseTrack:
			response, err := wru.websocketHandler.PauseTrack(context.Background(), room, user)
			wru.msgBrokerUtil.BroadcastMessage(room.ID, response)
			if err != nil {
				wru.log.Info(op, "error", err.Error(), "user", user.Username)
			}
		case domain.WSRoomSeekTrack:
			var req domain.WSRoomSeekTrackRequest
			req.Time = int64(message.Data.(map[string]interface{})["time"].(float64))
			response, err := wru.websocketHandler.SeekTrack(context.Background(), room, user, req)
			wru.msgBrokerUtil.BroadcastMessage(room.ID, response)
			if err != nil {
				wru.log.Info(op, "error", err.Error(), "user", user.Username)
			}
		case domain.WSRoomSyncTrack:
			response, err := wru.websocketHandler.SyncTrack(context.Background(), room, user)
			wru.msgBrokerUtil.BroadcastMessage(room.ID, response)
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
		case domain.WSRoomStopTrack:
			response, err := wru.websocketHandler.StopTrack(context.Background(), room, user)
			wru.msgBrokerUtil.BroadcastMessage(room.ID, response)
			if err != nil {
				wru.log.Info(op, "error", err.Error(), "user", user.Username)
			}
		}
		duration := time.Since(startTime).Seconds()
		wru.prom.WebsocketMessageHandlingDuration.WithLabelValues("ws/room?id="+strconv.Itoa(int(room.ID)), strconv.Itoa(message.Type), strconv.Itoa(int(user.ID))).Observe(duration)
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
	response, err := wru.websocketHandler.LoggedIn(ctx, roomID, userID)
	wru.msgBrokerUtil.BroadcastMessage(roomID, response)
	if err != nil {
		wru.log.Info("Websockets.Room", "error", err.Error(), "user", userID)
	}
	wru.prom.WebsocketConnectionsCount.Inc()
}

func (wru *wsRoomUsecase) LoggedOut(ctx context.Context, roomID int64, userID int64) {
	ctx, cancel := context.WithTimeout(ctx, wru.timeout)
	defer cancel()
	wru.websocketHandler.Remove(roomID, userID)
	response, err := wru.websocketHandler.LoggedOut(ctx, roomID, userID)
	wru.msgBrokerUtil.BroadcastMessage(roomID, response)
	if err != nil {
		wru.log.Info("Websockets.Room", "error", err.Error(), "user", userID)
	}
	wru.prom.WebsocketConnectionsCount.Dec()
}
