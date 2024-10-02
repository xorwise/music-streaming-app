package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xorwise/music-streaming-service/api/routes"
	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
	"github.com/xorwise/music-streaming-service/internal/repository"
	"github.com/xorwise/music-streaming-service/internal/utils"
	"github.com/xorwise/music-streaming-service/internal/utils/websockets"
)

func main() {
	// Config
	cfg := bootstrap.NewConfig()
	timeout := time.Duration(cfg.RequestTimeout) * time.Second

	// Logging
	log := setupLogger()

	// Database
	db := bootstrap.NewDatabaseConnection(cfg)
	bootstrap.MigrateDatabase(db, cfg)

	// Workers
	trackCh := make(chan domain.TrackStatus)
	errorCh := make(chan error)
	clients := make(domain.WSClients)
	startWorkers(db, clients, trackCh, errorCh)

	// Prometheus
	prom := bootstrap.NewPrometheus()
	prom.Init()

	// Routers
	mux := http.NewServeMux()
	routes.Setup(cfg, timeout, db, mux, log, clients, trackCh, errorCh, prom)
	mux.Handle("/metrics", promhttp.Handler())

	defer db.Close()

	log.Info("starting server at", slog.Int("port", cfg.Port))

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), mux)
	if err != nil {
		log.Error("failed to start server", err.Error())
	}
}

func setupLogger() *slog.Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}

func startWorkers(db *sql.DB, clients domain.WSClients, trackCh chan domain.TrackStatus, errorCh chan error) {
	tr := repository.NewTrackRepository(db, trackCh)
	tu := utils.NewTrackUtils(errorCh)
	go tr.RemoveOutdated(tu)
	wsh := websockets.NewWebsocketHandler(clients, trackCh)
	go wsh.HandleTrackEvent()
}
