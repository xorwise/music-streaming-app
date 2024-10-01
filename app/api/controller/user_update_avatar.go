package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type UserUpdateAvatarController struct {
	Usecase domain.UserUpdateAvatarUsecase
	Cfg     *bootstrap.Config
	Log     *slog.Logger
}

func (uc *UserUpdateAvatarController) Handle(w http.ResponseWriter, r *http.Request) {
	const op = "Controllers.UserUpdateAvatar"
	ctx := r.Context()
	user := ctx.Value("user").(*domain.User)

	var req domain.UserUpdateAvatarRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		uc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	path, err := uc.Usecase.SaveFile(ctx, req.Data, req.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		uc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	user.Avatar = path
	err = uc.Usecase.Update(ctx, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.ErrorResponse{Error: err.Error()})
		uc.Log.Info(op, "error", err.Error(), "user", user.Username)
		return
	}

	var response domain.UserMeResponse

	response.ID = user.ID
	response.Username = user.Username
	response.Avatar = user.Avatar

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	uc.Log.Info(op, "user", user.Username)
}
