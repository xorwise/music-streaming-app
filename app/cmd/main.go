package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/nats-io/nats.go"
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
	port := flag.Int("port", cfg.Port, "test")
	flag.Parse()

	// Logging
	log := setupLogger()

	// Database
	db := bootstrap.NewDatabaseConnection(cfg)
	bootstrap.MigrateDatabase(db, cfg)

	// Nats
	conn, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		panic(err)
	}

	// Workers
	trackCh := make(chan domain.TrackStatus)
	errorCh := make(chan error)
	clients := make(domain.WSClients)
	mbu := startWorkers(db, clients, trackCh, errorCh, conn)

	// Prometheus
	prom := bootstrap.NewPrometheus()
	prom.Init()

	// Routers
	mux := http.NewServeMux()
	routes.Setup(cfg, timeout, db, mux, log, clients, trackCh, errorCh, prom, mbu)
	mux.Handle("/metrics", promhttp.Handler())

	defer db.Close()

	log.Info("started server on", "port", *port)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), mux)
	if err != nil {
		log.Error("failed to start server", err.Error())
	}
}

func setupLogger() *slog.Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}

func startWorkers(db *sql.DB, clients domain.WSClients, trackCh chan domain.TrackStatus, errorCh chan error, conn *nats.Conn) domain.MessageBrokerUtils {
	tr := repository.NewTrackRepository(db, trackCh)
	tu := utils.NewTrackUtils(errorCh)
	go tr.RemoveOutdated(tu)
	wsh := websockets.NewWebsocketHandler(clients, trackCh)
	go wsh.HandleTrackEvent()
	broadcastCh := make(chan *domain.RoomBroadcastResponse)
	mbu := utils.NewNatsUtils(conn, broadcastCh, wsh)
	err := mbu.SubscribeToNats()
	if err != nil {
		panic(err)
	}
	mbu.HandleRoomClientRequests()
	go wsh.BroadcastMsg(broadcastCh)
	return mbu
}
