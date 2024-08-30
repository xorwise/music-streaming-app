package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/xorwise/music-streaming-service/internal/app"
	"github.com/xorwise/music-streaming-service/internal/config"
)

func main() {
	cfg := config.New()

	log := setupLogger()

	application := app.New(log, cfg.GRPC.Port, *cfg)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	log.Info("stopped gRPC server")
}

func setupLogger() *slog.Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
