package websockets

import (
	"context"
	"errors"
	"sync"

	"github.com/xorwise/music-streaming-service/internal/domain"
	"golang.org/x/net/websocket"
)

type currentTrackStatus struct {
	Track  *domain.Track `json:"track"`
	Time   int64         `json:"time"`
	Status trackStatus   `json:"status"`
}

type webSocketHandler struct {
	clients domain.WSClients
	mutexes map[int64]*sync.Mutex
	tracks  map[int64]*currentTrackStatus
	trackCh chan domain.TrackStatus
}

func NewWebsocketHandler(clients domain.WSClients, trackCh chan domain.TrackStatus) domain.WebSocketHandler {
	return &webSocketHandler{
		clients: clients,
		mutexes: make(map[int64]*sync.Mutex),
		tracks:  make(map[int64]*currentTrackStatus),
		trackCh: trackCh,
	}
}

func (wsh *webSocketHandler) Add(roomID int64, userID int64, conn *websocket.Conn) {
	clients, ok := wsh.clients[roomID]
	if !ok {
		clients = make(map[int64]*websocket.Conn)
		wsh.clients[roomID] = clients
		wsh.mutexes[roomID] = &sync.Mutex{}
	}
	wsh.mutexes[roomID].Lock()
	clients[userID] = conn
	wsh.mutexes[roomID].Unlock()
	websocket.JSON.Send(clients[userID], domain.WSRoomResponse{
		Type: domain.WSRoomLoggedInTrack,
		Data: struct {
			*currentTrackStatus
		}{
			wsh.tracks[roomID],
		},
	})
}

func (wsh *webSocketHandler) Remove(roomID int64, userID int64) {
	for roomID := range wsh.clients {
		wsh.mutexes[roomID].Lock()
		delete(wsh.clients[roomID], userID)
		wsh.mutexes[roomID].Unlock()
	}
}

func (wsh *webSocketHandler) LoggedIn(ctx context.Context, roomID int64, userID int64) error {
	var response domain.WSRoomResponse
	response.Type = domain.WSRoomLoggedIn
	response.Data = struct {
		Event  string `json:"event"`
		UserID int64  `json:"userID"`
	}{
		Event:  "connected",
		UserID: userID,
	}

	clients, ok := wsh.clients[roomID]
	if !ok {
		websocket.JSON.Send(wsh.clients[roomID][userID], domain.WSRoomResponse{
			Type:  domain.WSRoomError,
			Data:  "",
			Error: "internal server error",
		})
		return errors.New("internal server error")
	}

	for _, client := range clients {
		websocket.JSON.Send(client, response)
	}
	return nil
}

func (wrh *webSocketHandler) LoggedOut(ctx context.Context, roomID int64, userID int64) error {
	var response domain.WSRoomResponse
	response.Type = domain.WSRoomLoggedOut
	response.Data = struct {
		Event  string `json:"event"`
		UserID int64  `json:"userID"`
	}{
		Event:  "disconnected",
		UserID: userID,
	}

	clients, ok := wrh.clients[roomID]
	if !ok {
		websocket.JSON.Send(wrh.clients[roomID][userID], domain.WSRoomResponse{
			Type:  domain.WSRoomError,
			Data:  "",
			Error: "internal server error",
		})
		return errors.New("internal server error")
	}

	for _, client := range clients {
		websocket.JSON.Send(client, response)
	}
	return nil
}

func (wrh *webSocketHandler) GetOnlineUsers(ctx context.Context, roomID int64, userID int64) error {
	clients, ok := wrh.clients[roomID]
	if !ok {
		websocket.JSON.Send(wrh.clients[roomID][userID], domain.WSRoomResponse{
			Type:  domain.WSRoomError,
			Data:  "",
			Error: "internal server error",
		})
		return errors.New("internal server error")
	}

	var response domain.WSRoomResponse
	response.Type = domain.WSRoomGetOnlineUsers
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
	return nil
}
