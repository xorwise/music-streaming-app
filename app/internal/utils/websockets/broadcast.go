package websockets

import (
	"github.com/xorwise/music-streaming-service/internal/domain"
	"golang.org/x/net/websocket"
)

func (wsh *webSocketHandler) BroadcastMsg(broadcast chan *domain.RoomBroadcastResponse) {
	for {
		response := <-broadcast
		for _, client := range wsh.clients[response.RoomID] {
			websocket.JSON.Send(client, response.Response)
		}
	}
}
