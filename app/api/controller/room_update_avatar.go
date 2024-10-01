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

type RoomUpdateAvatarController struct {
	Usecase domain.RoomUpdateAvatarUsecase
	Cfg     *bootstrap.Config
	Log     *slog.Logger
}

func (rc *RoomUpdateAvatarController) Handle(w http.ResponseWriter, r *http.Request) {
	const op = "Controllers.RoomUpdateAvatar"
	ctx := r.Context()
	user := ctx.Value("user").(*domain.User)

	var req domain.UserUpdateAvatarRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	room, err := rc.Usecase.GetByID(ctx, int64(id))
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

	_, err = rc.Usecase.GetByUserIDandRoomID(ctx, user.ID, room.ID)
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

	if room.OwnerID != user.ID {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: "you are not the owner of this room"})
		rc.Log.Info(op, "error", "you are not the owner of this room", "user", user.Username)
		return
	}

	path, err := rc.Usecase.SaveFile(ctx, req.Data, req.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	room.Avatar = path
	err = rc.Usecase.Update(ctx, room)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	var response domain.RoomCreateResponse
	response.ID = room.ID
	response.Name = room.Name
	response.OwnerID = room.OwnerID
	response.Avatar = room.Avatar
	response.Code = room.Code
	response.CreatedAt = room.CreatedAt
	response.UpdatedAt = room.UpdatedAt

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	rc.Log.Info(op, "user", user.Username)
}
