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
	prom *bootstrap.Prometheus,
	mbu domain.MessageBrokerUtils,
) {
	// Media files
	SetupFileServer(cfg, db, mux, log)

	// Public routes
	NewUserCreateRoute(cfg, timeout, db, mux, log, prom)
	NewUserLoginRoute(cfg, timeout, db, mux, log, prom)

	// Protected routes
	NewUserMeRoute(cfg, timeout, db, mux, log, prom)
	NewUserUpdateAvatarRoute(cfg, timeout, db, mux, log, prom)

	NewRoomCreateRoute(cfg, timeout, db, mux, log, prom)
	NewRoomUsersRoute(cfg, timeout, db, mux, log, prom)
	NewRoomEnterRoute(cfg, timeout, db, mux, log, prom)
	NewRoomListByUserRoute(cfg, timeout, db, mux, log, prom)
	NewRoomLeaveRoute(cfg, timeout, db, mux, log, prom)
	NewRoomGetByIDRoute(cfg, timeout, db, mux, log, prom)
	NewRoomUpdateAvatarRoute(cfg, timeout, db, mux, log, prom)

	NewTrackAddRoute(cfg, timeout, db, mux, log, trackCh, errorCh, prom)
	NewTrackListByRoomRoute(cfg, timeout, db, mux, log, trackCh, prom)
	NewTrackDeleteRoute(cfg, timeout, db, mux, log, trackCh, prom)

	// Websocket routes
	NewWSRoomRoute(cfg, timeout, db, mux, log, clients, trackCh, prom, mbu)
}
