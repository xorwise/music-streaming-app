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

type TrackDeleteController struct {
	Usecase domain.TrackDeleteUsecase
	Cfg     *bootstrap.Config
	Log     *slog.Logger
}

func (tc *TrackDeleteController) Handle(w http.ResponseWriter, r *http.Request) {
	const op = "Controllers.TrackDelete"

	ctx := r.Context()
	user := ctx.Value("user").(*domain.User)
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		tc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	track, err := tc.Usecase.GetByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, domain.ErrTrackNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		tc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	_, err = tc.Usecase.GetByUserIDandRoomID(ctx, track.RoomID, user.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotUserInRoom) {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		tc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	err = tc.Usecase.Remove(ctx, track)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		tc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	err = tc.Usecase.RemoveFiles(ctx, track)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		tc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	tc.Log.Info(op, "user", user.Username)
}
