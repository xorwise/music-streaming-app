package routes

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
)

func Setup(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger) {
	NewUserCreateRoute(cfg, timeout, db, mux, log)
	NewUserLoginRoute(cfg, timeout, db, mux, log)
}
