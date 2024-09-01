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

func NewUserCreateRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger) {
	ur := repository.NewUserRepository(db)
	uc := controller.UserCreateController{
		Usecase: usecase.NewUserCreateUsecase(ur, timeout),
		Cfg:     cfg,
	}
	mw := middleware.NewLoggingMiddleware(log)

	mux.Handle("POST /users", mw.Handle(http.HandlerFunc(uc.Handle)))
}
