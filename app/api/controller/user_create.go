package controller

import (
	"encoding/json"
	"net/http"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
	"github.com/xorwise/music-streaming-service/internal/utils"
)

type UserCreateController struct {
	Usecase domain.UserCreateUsecase
	Cfg     *bootstrap.Config
}

func (uc *UserCreateController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctx = r.Context()

	var request domain.UserCreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	passHashByte, err := utils.HashPassword(request.Password)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	user := &domain.User{
		Username: request.Username,
		PassHash: string(passHashByte),
	}

	id, err := uc.Usecase.Create(ctx, user)
	if err != nil {
		// TODO: custom errors
		json.NewEncoder(w).Encode(err)
		return
	}
	response := domain.UserCreateResponse{
		ID: id,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
