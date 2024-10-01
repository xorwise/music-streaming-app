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
	Stopped trackStatus = "stopped"
)

func (wrh *webSocketHandler) PlayTrack(ctx context.Context, room *domain.Room, track *domain.Track, message domain.WSRoomPlayTrackRequest) error {
	wrh.tracks[room.ID] = &currentTrackStatus{
		Track:  track,
		Time:   message.Time,
		Status: Playing,
	}
	for id := range wrh.clients[room.ID] {
		websocket.JSON.Send(wrh.clients[room.ID][id], domain.WSRoomResponse{
			Type: domain.WSRoomPlayTrack,
			Data: struct {
				domain.WSRoomPlayTrackRequest
				Path string `json:"path"`
			}{
				message,
				track.Path,
			},
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

	wrh.tracks[room.ID].Status = Paused

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
	wrh.tracks[room.ID].Time = message.Time
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
				Status: trStatus.Status,
				Time:   trStatus.Time,
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

	if message.Time > wrh.tracks[room.ID].Time {
		wrh.tracks[room.ID].Time = message.Time
	}

	return nil
}

func (wrh *webSocketHandler) StopTrack(ctx context.Context, room *domain.Room, user *domain.User) error {
	_, ok := wrh.tracks[room.ID]
	if !ok {
		websocket.JSON.Send(wrh.clients[room.ID][user.ID], domain.WSRoomResponse{
			Type:  domain.WSRoomError,
			Data:  "",
			Error: "there is not track playing",
		})
		return errors.New("there is not track playing")
	}
	wrh.tracks[room.ID].Status = Stopped
	wrh.tracks[room.ID].Time = 0

	for id := range wrh.clients[room.ID] {
		websocket.JSON.Send(wrh.clients[room.ID][id], domain.WSRoomResponse{
			Type:  domain.WSRoomStopTrack,
			Data:  "",
			Error: "",
		})
	}
	return nil

}
