package domain

import (
	"sync"

	"golang.org/x/net/websocket"
)

type WSClients struct {
	RoomClients map[int64]map[int64]*websocket.Conn
	Mutexes     map[int64]*sync.Mutex
}

func NewWSClients() *WSClients {
	return &WSClients{
		RoomClients: make(map[int64]map[int64]*websocket.Conn),
		Mutexes:     make(map[int64]*sync.Mutex),
	}
}

func (wsc *WSClients) Add(roomID int64, userID int64, conn *websocket.Conn) {
	clients, ok := wsc.RoomClients[roomID]
	if !ok {
		clients = make(map[int64]*websocket.Conn)
		wsc.RoomClients[roomID] = clients
		wsc.Mutexes[roomID] = &sync.Mutex{}
	}
	wsc.Mutexes[roomID].Lock()
	defer wsc.Mutexes[roomID].Unlock()
	clients[userID] = conn
}

func (wsc *WSClients) Remove(roomID int64, userID int64) {
	for roomID := range wsc.RoomClients {
		wsc.Mutexes[roomID].Lock()
		delete(wsc.RoomClients[roomID], userID)
		wsc.Mutexes[roomID].Unlock()
	}
}
