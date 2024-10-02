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
	"github.com/xorwise/music-streaming-service/internal/utils"
)

func NewTrackAddRoute(
	cfg *bootstrap.Config,
	timeout time.Duration,
	db *sql.DB,
	mux *http.ServeMux,
	log *slog.Logger,
	trackCh chan domain.TrackStatus,
	errorCh chan error,
	prom *bootstrap.Prometheus,
) {
	tr := repository.NewTrackRepository(db, trackCh)
	rr := repository.NewRoomRepository(db)
	tu := utils.NewTrackUtils(make(chan error))
	uc := controller.TrackAddController{
		Usecase: usecase.NewTrackAddUsecase(tr, rr, tu, errorCh, timeout),
		Cfg:     cfg,
		Log:     log,
	}
	lmw := middleware.NewLoggingMiddleware(log)

	ur := repository.NewUserRepository(db)
	jmw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)

	mmw := middleware.NewMetricsMiddleware(prom)

	mux.Handle("POST /tracks", mmw.Handle(jmw.LoginRequired(lmw.Handle(http.HandlerFunc(uc.Handle)))))
}
