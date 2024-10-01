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
	"github.com/xorwise/music-streaming-service/internal/utils"
)

func NewRoomUpdateAvatarRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger) {
	rr := repository.NewRoomRepository(db)
	uu := utils.NewUserUtils(cfg.TokenTTL, cfg.JWTSecret)
	uc := controller.RoomUpdateAvatarController{
		Usecase: usecase.NewRoomUpdateAvatarUsecase(rr, uu, timeout),
		Cfg:     cfg,
		Log:     log,
	}
	lmw := middleware.NewLoggingMiddleware(log)

	ur := repository.NewUserRepository(db)
	jmw := middleware.NewJWTMiddleware(cfg.JWTSecret, ur)

	mux.Handle("PUT /rooms/{id}/avatar", jmw.LoginRequired(lmw.Handle(http.HandlerFunc(uc.Handle))))
}
