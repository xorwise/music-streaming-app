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

func NewRoomGetByIDRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger, prom *bootstrap.Prometheus) {
	rr := repository.NewRoomRepository(db)
	uc := controller.RoomGetByIDController{
		Usecase: usecase.NewRoomGetByIDUsecase(rr, timeout),
		Cfg:     cfg,
		Log:     log,
	}
	mw := middleware.NewLoggingMiddleware(log)
	ur := repository.NewUserRepository(db)
	jmw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)

	mmw := middleware.NewMetricsMiddleware(prom)
	mux.Handle("GET /rooms/{id}", mmw.Handle(jmw.LoginRequired(mw.Handle(http.HandlerFunc(uc.Handle)))))
}
