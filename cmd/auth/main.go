package main

import (
	"log/slog"
	"os"

	"github.com/xorwise/music-streaming-service/internal/app"
	"github.com/xorwise/music-streaming-service/internal/config"
)

func main() {
	cfg := config.New()

	log := setupLogger()

	application := app.New(log, cfg.GRPC.Port, *cfg)
	application.GRPCServer.MustRun()

}

func setupLogger() *slog.Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
