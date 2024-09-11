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

func NewUserLoginRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger) {
	ur := repository.NewUserRepository(db)
	uc := controller.UserLoginController{
		Usecase: usecase.NewUserLoginUsecase(ur, timeout),
		Cfg:     cfg,
		Log:     log,
	}
	mw := middleware.NewLoggingMiddleware(log)

	mux.Handle("POST /users/login", mw.Handle(http.HandlerFunc(uc.Handle)))
}
