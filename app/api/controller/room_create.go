package controller

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type RoomCreateController struct {
	Usecase domain.RoomCreateUsecase
	Cfg     *bootstrap.Config
	Log     *slog.Logger
}

func (rc *RoomCreateController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctx = r.Context()
	user := r.Context().Value("user").(*domain.User)

	const op = "Room.Create"

	var request domain.RoomCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}
	if request.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: "room name is required"})
		rc.Log.Info(op, "error", "room name is required", "user", user.Username)
		return
	}
	room := &domain.Room{
		Name:      request.Name,
		OwnerID:   user.ID,
		CreatedAt: time.Now().UTC(),
	}
	id, err := rc.Usecase.Create(ctx, room)
	if err != nil {
		if errors.Is(err, domain.ErrFieldRequired) {
			w.WriteHeader(http.StatusBadRequest)
		} else if errors.Is(err, domain.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	code := rc.Usecase.GenerateCode(ctx, id)
	err = rc.Usecase.SetCode(ctx, id, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}
	err = rc.Usecase.AddRoomUser(ctx, id, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		rc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	response := domain.RoomCreateResponse{
		ID:        id,
		Name:      room.Name,
		Code:      code,
		OwnerID:   user.ID,
		CreatedAt: room.CreatedAt,
		UpdatedAt: room.UpdatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	rc.Log.Info(op, "creating room", room.Name, "user", user.Username)
}
