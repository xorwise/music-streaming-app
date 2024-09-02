package controller

import (
	"encoding/json"
	"net/http"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type UserMeController struct {
	Usecase domain.UserMeUsecase
	Cfg     *bootstrap.Config
}

func (c *UserMeController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user_id").(int64)

	user, err := c.Usecase.GetByID(r.Context(), userID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		return
	}

	var response domain.UserMeResponse
	response.ID = user.ID
	response.Username = user.Username

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
