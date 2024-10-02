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

func NewRoomLeaveRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger, prom *bootstrap.Prometheus) {
	rr := repository.NewRoomRepository(db)
	uc := controller.RoomLeaveController{
		Usecase: usecase.NewRoomLeaveUsecase(rr, timeout),
		Cfg:     cfg,
		Log:     log,
	}
	ur := repository.NewUserRepository(db)
	lmw := middleware.NewLoggingMiddleware(log)
	jmw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)

	mmw := middleware.NewMetricsMiddleware(prom)
	mux.Handle("POST /rooms/{id}/leave", mmw.Handle(jmw.LoginRequired(lmw.Handle(http.HandlerFunc(uc.Handle)))))
}
