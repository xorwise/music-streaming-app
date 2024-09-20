package controller

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type TrackListByRoomController struct {
	Usecase domain.TrackListByRoomUsecase
	Cfg     *bootstrap.Config
	Log     *slog.Logger
}

func (tc *TrackListByRoomController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	const op = "Track.ListByRoom"

	user := ctx.Value("user").(*domain.User)

	roomID, err := strconv.Atoi(r.PathValue("roomID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		tc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	params := r.URL.Query()

	room, err := tc.Usecase.GetRoomByID(ctx, int64(roomID))
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

	tracks, err := tc.Usecase.ListByRoomID(ctx, room.ID, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		tc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	var response []domain.TrackListByRoomResponse

	for _, track := range tracks {
		response = append(response, domain.TrackListByRoomResponse{
			ID:      track.ID,
			Title:   track.Title,
			Artist:  track.Artist,
			RoomID:  track.RoomID,
			Path:    track.Path,
			IsReady: track.IsReady,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
