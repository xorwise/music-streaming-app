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
	"github.com/xorwise/music-streaming-service/internal/utils/websockets"
)

func NewWSRoomRoute(
	cfg *bootstrap.Config,
	timeout time.Duration,
	db *sql.DB,
	mux *http.ServeMux,
	log *slog.Logger,
	clients domain.WSClients,
	trackCh chan domain.TrackStatus,
	prom *bootstrap.Prometheus,
	mbu domain.MessageBrokerUtils,
) {
	wsh := websockets.NewWebsocketHandler(clients, trackCh)
	lmw := middleware.NewLoggingMiddleware(log)

	ur := repository.NewUserRepository(db)
	jmw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)

	rr := repository.NewRoomRepository(db)
	tr := repository.NewTrackRepository(db, trackCh)

	wsmw := middleware.NewWSMiddleware()
	wsc := controller.WSRoomController{
		Usecase: usecase.NewWSRoomUsecase(rr, tr, wsh, mbu, log, prom, timeout),
		Cfg:     cfg,
		Log:     log,
	}

	mux.Handle("/room", jmw.LoginRequired(lmw.Handle(wsmw.Handle(wsc.Handle))))
}
