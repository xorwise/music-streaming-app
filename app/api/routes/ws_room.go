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

func NewWSRoomRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger, clients *domain.WSClients) {
	lmw := middleware.NewLoggingMiddleware(log)

	ur := repository.NewUserRepository(db)
	jmw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)

	rr := repository.NewRoomRepository(db)

	wsmw := middleware.NewWSMiddleware()
	wsc := controller.WSRoomController{
		Usecase: usecase.NewWSRoomUsecase(rr, timeout, clients, log),
		Cfg:     cfg,
		Clients: clients,
		Log:     log,
	}

	mux.Handle("/room", jmw.LoginRequired(lmw.Handle(wsmw.Handle(wsc.Handle))))
}
