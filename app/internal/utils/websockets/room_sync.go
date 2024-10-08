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

func (wrh *webSocketHandler) PlayTrack(ctx context.Context, room *domain.Room, track *domain.Track, message domain.WSRoomPlayTrackRequest) (*domain.WSRoomResponse, error) {
	wrh.tracks[room.ID] = &currentTrackStatus{
		Track:  track,
		Time:   message.Time,
		Status: Playing,
	}
	response := domain.WSRoomResponse{
		Type: domain.WSRoomPlayTrack,
		Data: struct {
			domain.WSRoomPlayTrackRequest
			Path string `json:"path"`
		}{
			message,
			track.Path,
		},
		Error: "",
	}
	return &response, nil
}

func (wrh *webSocketHandler) PauseTrack(ctx context.Context, room *domain.Room, user *domain.User) (*domain.WSRoomResponse, error) {
	_, ok := wrh.tracks[room.ID]
	if !ok {
		websocket.JSON.Send(wrh.clients[room.ID][user.ID], domain.WSRoomResponse{
			Type:  domain.WSRoomError,
			Data:  "",
			Error: "there is not track playing",
		})
		return nil, errors.New("there is not track playing")
	}

	wrh.tracks[room.ID].Status = Paused

	response := domain.WSRoomResponse{
		Type:  domain.WSRoomPauseTrack,
		Data:  "",
		Error: "",
	}

	return &response, nil
}

func (wrh *webSocketHandler) SeekTrack(ctx context.Context, room *domain.Room, user *domain.User, message domain.WSRoomSeekTrackRequest) (*domain.WSRoomResponse, error) {
	wrh.tracks[room.ID].Time = message.Time
	response := domain.WSRoomResponse{
		Type:  domain.WSRoomSeekTrack,
		Data:  message,
		Error: "",
	}

	return &response, nil
}

func (wrh *webSocketHandler) SyncTrack(ctx context.Context, room *domain.Room, user *domain.User) (*domain.WSRoomResponse, error) {
	trStatus, ok := wrh.tracks[room.ID]
	if !ok {
		websocket.JSON.Send(wrh.clients[room.ID][user.ID], domain.WSRoomResponse{
			Type:  domain.WSRoomError,
			Data:  "",
			Error: "there is not track playing",
		})
		return nil, errors.New("there is not track playing")
	}

	response := domain.WSRoomResponse{
		Type: domain.WSRoomSyncTrack,
		Data: struct {
			Status trackStatus `json:"status"`
			Time   int64       `json:"time"`
		}{
			Status: trStatus.Status,
			Time:   trStatus.Time,
		},
		Error: "",
	}

	return &response, nil
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

func (wrh *webSocketHandler) StopTrack(ctx context.Context, room *domain.Room, user *domain.User) (*domain.WSRoomResponse, error) {
	_, ok := wrh.tracks[room.ID]
	if !ok {
		websocket.JSON.Send(wrh.clients[room.ID][user.ID], domain.WSRoomResponse{
			Type:  domain.WSRoomError,
			Data:  "",
			Error: "there is not track playing",
		})
		return nil, errors.New("there is not track playing")
	}
	wrh.tracks[room.ID].Status = Stopped
	wrh.tracks[room.ID].Time = 0

	response := domain.WSRoomResponse{
		Type:  domain.WSRoomStopTrack,
		Data:  "",
		Error: "",
	}

	return &response, nil
}
