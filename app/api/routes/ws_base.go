package routes

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/xorwise/music-streaming-service/api/controller"
	"github.com/xorwise/music-streaming-service/api/middleware"
	"github.com/xorwise/music-streaming-service/internal/bootstrap"
)

func NewWSBaseRoute(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger) {
	lmw := middleware.NewLoggingMiddleware(log)
	wsMw := middleware.NewWSMiddleware()
	wsc := controller.WSBaseController{
		Cfg: cfg,
	}

	mux.Handle("/ws", lmw.Handle(wsMw.Handle(wsc.Handle)))
}
