package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
	"github.com/xorwise/music-streaming-service/internal/utils"
)

type UserLoginController struct {
	Usecase domain.UserLoginUsecase
	Cfg     *bootstrap.Config
	Log     *slog.Logger
}

func (uc *UserLoginController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctx = r.Context()

	const op = "User.Login"

	var request domain.UserLoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		uc.Log.Info(op, slog.With("error", err.Error()))
		return
	}

	user, err := uc.Usecase.GetByUsername(ctx, request.Username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		uc.Log.Info(op, slog.With("error", err.Error()))
		return
	}

	if ok := utils.CheckPasswordHash(request.Password, user.PassHash); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: "invalid credentials"})
		uc.Log.Info(op, slog.With("error", "invalid credentials"))
		return
	}

	tokenStr, err := uc.Usecase.CreateAccessToken(ctx, uc.Cfg, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		uc.Log.Info(op, slog.With("error", err.Error()))
		return
	}
	response := domain.UserLoginResponse{
		AccessToken: tokenStr,
	}

	uc.Log.Info(op, "user logged in, username", user.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
