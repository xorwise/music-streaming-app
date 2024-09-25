package websockets

import (
	"context"
	"errors"

	"github.com/xorwise/music-streaming-service/internal/domain"
	"golang.org/x/net/websocket"
)

type trackStatus string

const (
	Playing trackStatus = "playing"
	Paused  trackStatus = "paused"
	NotPlay trackStatus = "notplay"
)

func (wrh *webSocketHandler) PlayTrack(ctx context.Context, room *domain.Room, track *domain.Track, message domain.WSRoomPlayTrackRequest) error {
	wrh.tracks[room.ID] = &currentTrackStatus{
		track:  track,
		time:   message.Time,
		status: Playing,
	}
	for id := range wrh.clients[room.ID] {
		websocket.JSON.Send(wrh.clients[room.ID][id], domain.WSRoomResponse{
			Type:  domain.WSRoomPlayTrack,
			Data:  message,
			Error: "",
		})
	}
	return nil
}

func (wrh *webSocketHandler) PauseTrack(ctx context.Context, room *domain.Room, user *domain.User) error {
	_, ok := wrh.tracks[room.ID]
	if !ok {
		websocket.JSON.Send(wrh.clients[room.ID][user.ID], domain.WSRoomResponse{
			Type:  domain.WSRoomError,
			Data:  "",
			Error: "there is not track playing",
		})
		return errors.New("there is not track playing")
	}

	wrh.tracks[room.ID].status = Paused

	for id := range wrh.clients[room.ID] {
		websocket.JSON.Send(wrh.clients[room.ID][id], domain.WSRoomResponse{
			Type:  domain.WSRoomPauseTrack,
			Data:  "",
			Error: "",
		})
	}
	return nil

}

func (wrh *webSocketHandler) SeekTrack(ctx context.Context, room *domain.Room, user *domain.User, message domain.WSRoomSeekTrackRequest) error {
	wrh.tracks[room.ID].time = message.Time
	for id := range wrh.clients[room.ID] {
		websocket.JSON.Send(wrh.clients[room.ID][id], domain.WSRoomResponse{
			Type:  domain.WSRoomSeekTrack,
			Data:  message,
			Error: "",
		})
	}
	return nil
}

func (wrh *webSocketHandler) SyncTrack(ctx context.Context, room *domain.Room, user *domain.User) error {
	trStatus, ok := wrh.tracks[room.ID]
	if !ok {
		websocket.JSON.Send(wrh.clients[room.ID][user.ID], domain.WSRoomResponse{
			Type:  domain.WSRoomError,
			Data:  "",
			Error: "there is not track playing",
		})
		return errors.New("there is not track playing")
	}
	for id := range wrh.clients[room.ID] {
		websocket.JSON.Send(wrh.clients[room.ID][id], domain.WSRoomResponse{
			Type: domain.WSRoomSyncTrack,
			Data: struct {
				Status trackStatus `json:"status"`
				Time   int64       `json:"time"`
			}{
				Status: trStatus.status,
				Time:   trStatus.time,
			},
			Error: "",
		})
	}
	return nil
}

func (wrh *webSocketHandler) UpdateTrackTime(ctx context.Context, room *domain.Room, user *domain.User, message domain.WSRoomUpdateTrackTimeRequest) error {
	_, ok := wrh.tracks[room.ID]
	if !ok {
		websocket.JSON.Send(wrh.clients[room.ID][user.ID], domain.WSRoomResponse{
			Type:  domain.WSRoomError,
			Data:  "",
			Error: "there is not track playing",
		})
		return errors.New("there is not track playing")
	}

	wrh.tracks[room.ID].time = message.Time

	return nil
}
