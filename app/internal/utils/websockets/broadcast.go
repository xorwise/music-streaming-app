package websockets

import (
	"context"
	"database/sql"

	"github.com/xorwise/music-streaming-service/internal/domain"
	"github.com/xorwise/music-streaming-service/internal/repository"
	"golang.org/x/net/websocket"
)

func (wsh *webSocketHandler) BroadcastMsg(broadcast chan *domain.RoomBroadcastResponse, db *sql.DB) {
	trackCh := make(chan domain.TrackStatus)
	repo := repository.NewTrackRepository(db, trackCh)
	for {
		response := <-broadcast
		switch response.Response.Type {
		case domain.WSRoomPlayTrack:
			trackID := int64(response.Response.Data.(map[string]interface{})["trackID"].(float64))
			time := int64(response.Response.Data.(map[string]interface{})["time"].(float64))
			track, err := repo.GetByID(context.Background(), trackID)
			if err != nil {
				break
			}
			wsh.tracks[response.RoomID] = &currentTrackStatus{
				Track:  track,
				Time:   time,
				Status: Playing,
			}
		case domain.WSRoomPauseTrack:
			_, ok := wsh.tracks[response.RoomID]
			if !ok {
				break
			}
			wsh.tracks[response.RoomID].Status = Paused
		case domain.WSRoomSeekTrack:
			time := int64(response.Response.Data.(map[string]interface{})["time"].(float64))
			_, ok := wsh.tracks[response.RoomID]
			if !ok {
				break
			}
			wsh.tracks[response.RoomID].Time = time
		case domain.WSRoomUpdateTrackTime:
			time := int64(response.Response.Data.(map[string]interface{})["time"].(float64))
			_, ok := wsh.tracks[response.RoomID]
			if !ok {
				break
			}
			if time > wsh.tracks[response.RoomID].Time {
				wsh.tracks[response.RoomID].Time = time
			}
		case domain.WSRoomStopTrack:
			_, ok := wsh.tracks[response.RoomID]
			if !ok {
				break
			}
			wsh.tracks[response.RoomID].Status = Stopped
			wsh.tracks[response.RoomID].Time = 0
		}
		for _, client := range wsh.clients[response.RoomID] {
			websocket.JSON.Send(client, response.Response)
		}
	}
}
