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

type RoomLeaveController struct {
	Usecase domain.RoomLeaveUsecase
	Cfg     *bootstrap.Config
	Log     *slog.Logger
}

func (rc *RoomLeaveController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ctx = r.Context()
	user := r.Context().Value("user").(*domain.User)

	const op = "Room.Leave"

	roomID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}
	id := int64(roomID)

	room, err := rc.Usecase.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrRoomNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	_, err = rc.Usecase.GetUserIDandRoomID(ctx, room.ID, user.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotUserInRoom) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	err = rc.Usecase.RemoveRoomUser(ctx, room.ID, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	w.WriteHeader(http.StatusOK)
	rc.Log.Info(op, "room", room.Name, "user", user.Username)
}
