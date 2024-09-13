package controller

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type TrackAddController struct {
	Usecase domain.TrackAddUsecase
	Cfg     *bootstrap.Config
	Log     *slog.Logger
}

func (tc *TrackAddController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctx = r.Context()

	const op = "Track.Add"

	user := ctx.Value("user").(*domain.User)

	var request domain.TrackAddRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		tc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	room, err := tc.Usecase.GetRoomByID(ctx, request.RoomID)
	if err != nil {
		if errors.Is(err, domain.ErrRoomNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		tc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	_, err = tc.Usecase.GetByUserIDandRoomID(ctx, room.ID, user.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotUserInRoom) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		tc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	trackCh := make(chan error, 1)
	path, err := tc.Usecase.FindAndSaveTrack(ctx, trackCh, request.Title, request.Artist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		tc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	track := &domain.Track{
		Title:  request.Title,
		Artist: request.Artist,
		Path:   path,
		RoomID: room.ID,
	}
	id, err := tc.Usecase.Create(ctx, track)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		tc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	go tc.Usecase.WaitForTrack(ctx, trackCh, track)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(domain.TrackAddResponse{ID: id})
}
