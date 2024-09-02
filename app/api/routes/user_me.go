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

func NewUserMeRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger) {
	ur := repository.NewUserRepository(db)
	uc := controller.UserMeController{
		Usecase: usecase.NewUserMeUsecase(ur, timeout),
		Cfg:     cfg,
	}
	loggingMw := middleware.NewLoggingMiddleware(log)
	jwtMw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)

	mux.Handle("GET /users/me", loggingMw.Handle(jwtMw.LoginRequired(http.HandlerFunc(uc.Handle))))
}
