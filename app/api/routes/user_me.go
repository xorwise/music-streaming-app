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
)

func NewUserMeRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger, prom *bootstrap.Prometheus) {
	ur := repository.NewUserRepository(db)
	uc := controller.UserMeController{
		Usecase: usecase.NewUserMeUsecase(ur, timeout),
		Cfg:     cfg,
		Log:     log,
	}
	loggingMw := middleware.NewLoggingMiddleware(log)
	jwtMw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)
	mmw := middleware.NewMetricsMiddleware(prom)

	mux.Handle("GET /users/me", mmw.Handle(jwtMw.LoginRequired(loggingMw.Handle(http.HandlerFunc(uc.Handle)))))
}
