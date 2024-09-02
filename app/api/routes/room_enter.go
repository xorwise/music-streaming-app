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

func NewRoomEnterRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger) {
	rr := repository.NewRoomRepository(db)
	uc := controller.RoomEnterController{
		Usecase: usecase.NewRoomEnterUsecase(rr, timeout),
		Cfg:     cfg,
	}
	ur := repository.NewUserRepository(db)
	lmw := middleware.NewLoggingMiddleware(log)
	jmw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)

	mux.Handle("POST /rooms/enter", jmw.LoginRequired(lmw.Handle(http.HandlerFunc(uc.Handle))))
}
