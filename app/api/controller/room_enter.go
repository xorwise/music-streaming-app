package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type RoomEnterController struct {
	Usecase domain.RoomEnterUsecase
	Cfg     *bootstrap.Config
}

func (rc *RoomEnterController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctx = r.Context()

	user := r.Context().Value("user").(*domain.User)
	var request domain.RoomEnterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		return
	}
	code := request.Code

	room, err := rc.Usecase.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, domain.ErrRoomNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		return
	}

	_, err = rc.Usecase.GetByUserIDandRoomID(ctx, room.ID, user.ID)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: "you are already in this room"})
		return
	}

	err = rc.Usecase.AddRoomUser(ctx, room.ID, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
