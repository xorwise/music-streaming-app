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

func NewUserLoginRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger, prom *bootstrap.Prometheus) {
	ur := repository.NewUserRepository(db)
	uu := utils.NewUserUtils(cfg.TokenTTL, cfg.JWTSecret)
	uc := controller.UserLoginController{
		Usecase: usecase.NewUserLoginUsecase(ur, uu, timeout),
		Cfg:     cfg,
		Log:     log,
	}
	lmw := middleware.NewLoggingMiddleware(log)
	mmw := middleware.NewMetricsMiddleware(prom)

	mux.Handle("POST /users/login", mmw.Handle(lmw.Handle(http.HandlerFunc(uc.Handle))))
}
