package controller

import (
	"encoding/json"
	"net/http"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
	"github.com/xorwise/music-streaming-service/internal/utils"
)

type UserLoginController struct {
	Usecase domain.UserLoginUsecase
	Cfg     *bootstrap.Config
}

func (uc *UserLoginController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctx = r.Context()

	var request domain.UserLoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	user, err := uc.Usecase.GetByUsername(ctx, request.Username)
	if err != nil {
		// TODO: custom errors
		json.NewEncoder(w).Encode(err)
		return
	}

	if ok := utils.CheckPasswordHash(request.Password, user.PassHash); !ok {
		json.NewEncoder(w).Encode("invalid credentials")
		return
	}

	tokenStr, err := uc.Usecase.CreateAccessToken(ctx, uc.Cfg, user)
	if err != nil {
		// TODO: custom errors
		json.NewEncoder(w).Encode(err)
		return
	}
	response := domain.UserLoginResponse{
		AccessToken: tokenStr,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
