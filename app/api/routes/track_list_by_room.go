package routes

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/xorwise/music-streaming-service/api/controller"
	"github.com/xorwise/music-streaming-service/api/middleware"
	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
	"github.com/xorwise/music-streaming-service/internal/repository"
	"github.com/xorwise/music-streaming-service/internal/usecase"
)

func NewTrackListByRoomRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger, trackCh chan domain.TrackStatus) {
	tr := repository.NewTrackRepository(db, trackCh)
	rr := repository.NewRoomRepository(db)
	uc := controller.TrackListByRoomController{
		Usecase: usecase.NewTrackListByRoomUsecase(tr, rr, timeout),
		Cfg:     cfg,
		Log:     log,
	}
	lmw := middleware.NewLoggingMiddleware(log)

	ur := repository.NewUserRepository(db)
	jmw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)

	mux.Handle("GET /tracks/room/{roomID}", jmw.LoginRequired(lmw.Handle(http.HandlerFunc(uc.Handle))))
}
