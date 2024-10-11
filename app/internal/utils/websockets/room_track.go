package websockets

import (
	"github.com/xorwise/music-streaming-service/internal/domain"
	"golang.org/x/net/websocket"
)

func (wrh *webSocketHandler) HandleTrackEvent() {
	for {
		trackEvent := <-wrh.trackCh
		var event string
		if trackEvent.IsReady {
			event = "ready"
		} else if trackEvent.Path == "" {
			event = "expired"
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
}
