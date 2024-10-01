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

type RoomGetByIDController struct {
	Usecase domain.RoomGetByIDUsecase
	Cfg     *bootstrap.Config
	Log     *slog.Logger
}

func (rc *RoomGetByIDController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	user := ctx.Value("user").(*domain.User)

	const op = "Room.GetByID"

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

	var response domain.RoomGetByIDResponse
	response.ID = room.ID
	response.Name = room.Name
	response.Avatar = room.Avatar
	response.Code = room.Code
	response.OwnerID = room.OwnerID
	response.CreatedAt = room.CreatedAt
	response.UpdatedAt = room.UpdatedAt
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	rc.Log.Info(op, "room", room.ID, "user", user.Username)

}
