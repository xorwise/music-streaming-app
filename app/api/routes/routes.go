package routes

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
	"github.com/xorwise/music-streaming-service/internal/utils/websockets"
)

func Setup(cfg *bootstrap.Config, timeout time.Duration, db *sql.DB, mux *http.ServeMux, log *slog.Logger) {
	SetupFileServer(cfg, db, mux, log)
	NewUserCreateRoute(cfg, timeout, db, mux, log)
	NewUserLoginRoute(cfg, timeout, db, mux, log)

	// Protected routes
	NewUserMeRoute(cfg, timeout, db, mux, log)

	NewRoomCreateRoute(cfg, timeout, db, mux, log)
	NewRoomUsersRoute(cfg, timeout, db, mux, log)
	NewRoomEnterRoute(cfg, timeout, db, mux, log)
	NewRoomListByUserRoute(cfg, timeout, db, mux, log)
	NewRoomLeaveRoute(cfg, timeout, db, mux, log)
	NewRoomGetByIDRoute(cfg, timeout, db, mux, log)

	trackCh := make(chan domain.TrackStatus)
	NewTrackAddRoute(cfg, timeout, db, mux, log, trackCh)
	NewTrackListByRoomRoute(cfg, timeout, db, mux, log, trackCh)

	// Websocket routes
	wsh := websockets.NewWebsocketHandler(trackCh)
	go wsh.HandleTrackEvent()
	NewWSRoomRoute(cfg, timeout, db, mux, log, wsh, trackCh)
}
