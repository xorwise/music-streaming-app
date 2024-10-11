package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type UserMeController struct {
	Usecase domain.UserMeUsecase
	Cfg     *bootstrap.Config
	Log     *slog.Logger
}

func (c *UserMeController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user_id").(int64)

	const op = "User.Me"

	user, err := c.Usecase.GetByID(r.Context(), userID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		c.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	var response domain.UserMeResponse
	response.ID = user.ID
	response.Username = user.Username
	response.Avatar = user.Avatar

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	c.Log.Info(op, "user", user.Username)
}
