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

func NewTrackDeleteRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger, trackCh chan domain.TrackStatus, prom *bootstrap.Prometheus) {
	lmw := middleware.NewLoggingMiddleware(log)

	ur := repository.NewUserRepository(db)
	jmw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)

	tu := utils.NewTrackUtils(make(chan error))
	tr := repository.NewTrackRepository(db, trackCh)
	rr := repository.NewRoomRepository(db)
	tc := controller.TrackDeleteController{
		Usecase: usecase.NewTrackDeleteUsecase(tr, rr, tu, timeout),
		Cfg:     cfg,
		Log:     log,
	}
	mmw := middleware.NewMetricsMiddleware(prom)
	mux.Handle("DELETE /tracks/delete/{id}", mmw.Handle(jmw.LoginRequired(lmw.Handle(http.HandlerFunc(tc.Handle)))))
}
