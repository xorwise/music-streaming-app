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

type RoomUsersController struct {
	Usecase domain.RoomUsersUsecase
	Cfg     *bootstrap.Config
	Log     *slog.Logger
}

func (rc *RoomUsersController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctx = r.Context()
	user := r.Context().Value("user").(*domain.User)

	const op = "Room.Users"

	roomID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}
	id := int64(roomID)
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	_, err = rc.Usecase.GetByID(ctx, id)
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

	_, err = rc.Usecase.GetByUserIDandRoomID(ctx, id, user.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotUserInRoom) {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	users, err := rc.Usecase.ListRoomUsers(ctx, id, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}
	var response domain.RoomUsersResponse
	for _, user := range users {
		response.Users = append(response.Users, domain.UserMeResponse{
			ID:       user.ID,
			Username: user.Username,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	rc.Log.Info(op, "room", id, "user", user.Username)
}
