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

func NewUserUpdateAvatarRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger, prom *bootstrap.Prometheus) {
	ur := repository.NewUserRepository(db)
	jmw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)
	uu := utils.NewUserUtils(cfg.TokenTTL, cfg.JWTSecret)
	uc := controller.UserUpdateAvatarController{
		Usecase: usecase.NewUserUpdateAvatarUsecase(ur, uu, timeout),
		Cfg:     cfg,
		Log:     log,
	}
	lmw := middleware.NewLoggingMiddleware(log)
	mmw := middleware.NewMetricsMiddleware(prom)

	mux.Handle("PUT /users/avatar", mmw.Handle(jmw.LoginRequired(lmw.Handle(http.HandlerFunc(uc.Handle)))))
}
