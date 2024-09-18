package routes

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/xorwise/music-streaming-service/api/middleware"
	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/repository"
)

func SetupFileServer(cfg *bootstrap.Config, db *sql.DB, mux *http.ServeMux, log *slog.Logger) {
	ur := repository.NewUserRepository(db)
	_ = middleware.NewJWTMiddleware(cfg.JWTSecret, ur)

	lmw := middleware.NewLoggingMiddleware(log)

	fs := http.FileServer(http.Dir("./media"))

	mux.Handle("/media/", lmw.Handle(http.StripPrefix("/media/", fs)))
}
