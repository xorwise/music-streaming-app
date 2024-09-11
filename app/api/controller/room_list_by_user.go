package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type RoomListByUserController struct {
	Usecase domain.RoomListByUserUsecase
	Cfg     *bootstrap.Config
	Log     *slog.Logger
}

func (uc RoomListByUserController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	user := r.Context().Value("user").(*domain.User)

	const op = "Room.ListByUser"

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	userID := r.Context().Value("user_id").(int64)
	rooms, err := uc.Usecase.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		uc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	var response domain.RoomListByUserResponse
	for _, room := range rooms {
		response.Rooms = append(response.Rooms, &domain.RoomCreateResponse{
			ID:        room.ID,
			Name:      room.Name,
			Code:      room.Code,
			CreatedAt: room.CreatedAt,
			UpdatedAt: room.UpdatedAt,
			OwnerID:   room.OwnerID,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	uc.Log.Info(op, "user", user.Username)

}
