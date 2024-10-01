package controller

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type UserCreateController struct {
	Usecase domain.UserCreateUsecase
	Cfg     *bootstrap.Config
	Log     *slog.Logger
}

func (uc *UserCreateController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctx = r.Context()

	const op = "User.Create"

	var request domain.UserCreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		uc.Log.Info(op, slog.With("error", err.Error()))
		return
	}

	passHashByte, err := uc.Usecase.HashPassword(request.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		uc.Log.Info(op, slog.With("error", err.Error()))
		return
	}

	user := &domain.User{
		Username: request.Username,
		PassHash: string(passHashByte),
	}

	id, err := uc.Usecase.Create(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
		} else if errors.Is(err, domain.ErrFieldRequired) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		uc.Log.Info(op, slog.With("error", err.Error()))
		return
	}
	response := domain.UserCreateResponse{
		ID: id,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	uc.Log.Info(op, "creating user", user.Username)
}
