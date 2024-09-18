package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/xorwise/music-streaming-service/api/routes"
	"github.com/xorwise/music-streaming-service/internal/bootstrap"
)

func main() {
	cfg := bootstrap.NewConfig()

	log := setupLogger()

	db := bootstrap.NewDatabaseConnection(cfg)
	bootstrap.MigrateDatabase(db, cfg)
	timeout := time.Duration(cfg.RequestTimeout) * time.Second

	mux := http.NewServeMux()

	routes.Setup(cfg, timeout, db, mux, log)

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
