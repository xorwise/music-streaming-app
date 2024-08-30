package app

import (
	"log/slog"

	grpcapp "github.com/xorwise/music-streaming-service/internal/app/auth/grpc"
	"github.com/xorwise/music-streaming-service/internal/config/auth"
	"github.com/xorwise/music-streaming-service/internal/database/auth/postgresql"
	"github.com/xorwise/music-streaming-service/internal/services/auth"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, cfg config.Config) *App {
	database, err := postgresql.New(cfg)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, database, cfg.TokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}