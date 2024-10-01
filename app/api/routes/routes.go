package routes

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

func Setup(
	cfg *bootstrap.Config,
	timeout time.Duration,
	db *sql.DB, mux *http.ServeMux,
	log *slog.Logger,
	clients domain.WSClients,
	trackCh chan domain.TrackStatus,
	errorCh chan error,
) {
	// Media files
	SetupFileServer(cfg, db, mux, log)

	// Public routes
	NewUserCreateRoute(cfg, timeout, db, mux, log)
	NewUserLoginRoute(cfg, timeout, db, mux, log)

	// Protected routes
	NewUserMeRoute(cfg, timeout, db, mux, log)
	NewUserUpdateAvatarRoute(cfg, timeout, db, mux, log)

	NewRoomCreateRoute(cfg, timeout, db, mux, log)
	NewRoomUsersRoute(cfg, timeout, db, mux, log)
	NewRoomEnterRoute(cfg, timeout, db, mux, log)
	NewRoomListByUserRoute(cfg, timeout, db, mux, log)
	NewRoomLeaveRoute(cfg, timeout, db, mux, log)
	NewRoomGetByIDRoute(cfg, timeout, db, mux, log)
	NewRoomUpdateAvatarRoute(cfg, timeout, db, mux, log)

	NewTrackAddRoute(cfg, timeout, db, mux, log, trackCh, errorCh)
	NewTrackListByRoomRoute(cfg, timeout, db, mux, log, trackCh)
	NewTrackDeleteRoute(cfg, timeout, db, mux, log, trackCh)

	// Websocket routes
	NewWSRoomRoute(cfg, timeout, db, mux, log, clients, trackCh)
}
