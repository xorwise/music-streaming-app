package routes

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/xorwise/music-streaming-service/api/controller"
	"github.com/xorwise/music-streaming-service/api/middleware"
	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/repository"
	"github.com/xorwise/music-streaming-service/internal/usecase"
	"github.com/xorwise/music-streaming-service/internal/utils"
)

func NewUserUpdateAvatarRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger) {
	ur := repository.NewUserRepository(db)
	jmw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)
	uu := utils.NewUserUtils(cfg.TokenTTL, cfg.JWTSecret)
	uc := controller.UserUpdateAvatarController{
		Usecase: usecase.NewUserUpdateAvatarUsecase(ur, uu, timeout),
		Cfg:     cfg,
		Log:     log,
	}
	mw := middleware.NewLoggingMiddleware(log)

	mux.Handle("PUT /users/avatar", jmw.LoginRequired(mw.Handle(http.HandlerFunc(uc.Handle))))
}
