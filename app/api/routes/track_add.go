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

func NewTrackAddRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger) {
	tr := repository.NewTrackRepository(db)
	rr := repository.NewRoomRepository(db)
	uc := controller.TrackAddController{
		Usecase: usecase.NewTrackAddUsecase(tr, rr, timeout),
		Cfg:     cfg,
		Log:     log,
	}
	mw := middleware.NewLoggingMiddleware(log)

	ur := repository.NewUserRepository(db)
	jmw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)

	mux.Handle("POST /tracks", jmw.LoginRequired(mw.Handle(http.HandlerFunc(uc.Handle))))
}
