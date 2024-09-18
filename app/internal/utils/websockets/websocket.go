package websockets

import (
	"context"
	"errors"
	"sync"

	"github.com/xorwise/music-streaming-service/internal/domain"
	"golang.org/x/net/websocket"
)

const BufferSize = 6 * 1024

type webSocketHandler struct {
	clients map[int64]map[int64]*websocket.Conn
	mutexes map[int64]*sync.Mutex
	trackCh chan domain.TrackStatus
}

func NewWebsocketHandler(trackCh chan domain.TrackStatus) domain.WebSocketHandler {
	return &webSocketHandler{
		clients: make(map[int64]map[int64]*websocket.Conn),
		mutexes: make(map[int64]*sync.Mutex),
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
	defer wsh.mutexes[roomID].Unlock()
	clients[userID] = conn
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

// func (wrh *webSocketHandler) FetchMusicChunks(ctx context.Context, track *domain.Track, roomID int64, userID int64) error {
// 	musicReader, err := utils.NewMusicReader(track.Path)
// 	if err != nil {
// 		websocket.JSON.Send(wrh.clients[roomID][userID], domain.WSRoomResponse{
// 			Type:  domain.WSRoomError,
// 			Data:  "",
// 			Error: err.Error(),
// 		})
// 		return err
// 	}
// 	defer musicReader.Close()
//
// 	buff := make([]byte, BufferSize)
//
// 	for {
// 		n, err := musicReader.Read(buff)
// 		if err == io.EOF {
// 			return nil
// 		}
// 		if err != nil {
// 			return err
// 		}
//
// 		websocket.Message.Send(wrh.clients[roomID][userID], buff[:n])
// 	}
// }

func (wrh *webSocketHandler) HandleTrackEvent() {
	trackEvent := <-wrh.trackCh
	var event string
	if trackEvent.IsReady {
		event = "ready"

	} else {
		event = "removed"
	}
	msg := domain.WSRoomResponse{
		Type: domain.WSRoomTrackEvent,
		Data: struct {
			TrackID int64  `json:"trackID"`
			Path    string `json:"path"`
			Event   string `json:"event"`
		}{
			TrackID: trackEvent.ID,
			Path:    trackEvent.Path,
			Event:   event,
		},
	}
	for _, client := range wrh.clients[trackEvent.RoomID] {
		websocket.JSON.Send(client, msg)
	}
}
