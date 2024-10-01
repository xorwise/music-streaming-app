package routes

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/xorwise/music-streaming-service/api/middleware"
	"github.com/xorwise/music-streaming-service/internal/bootstrap"
)

func SetupFileServer(cfg *bootstrap.Config, db *sql.DB, mux *http.ServeMux, log *slog.Logger) {
	lmw := middleware.NewLoggingMiddleware(log)

	fs := http.FileServer(http.Dir("./media"))

	mux.Handle("/media/", lmw.Handle(http.StripPrefix("/media/", fs)))
}
