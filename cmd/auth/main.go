package main

import (
	"log/slog"
	"os"

	"github.com/xorwise/music-streaming-service/internal/app/auth"
	"github.com/xorwise/music-streaming-service/internal/config/auth"
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
