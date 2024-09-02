package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type RoomCreateController struct {
	Usecase domain.RoomCreateUsecase
	Cfg     *bootstrap.Config
}

func (rc *RoomCreateController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctx = r.Context()
	user := r.Context().Value("user").(*domain.User)

	var request domain.RoomCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
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
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	code := rc.Usecase.GenerateCode(ctx, id)
	err = rc.Usecase.SetCode(ctx, id, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	fmt.Println(user.ID, id)
	err = rc.Usecase.AddRoomUser(ctx, id, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
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
}
