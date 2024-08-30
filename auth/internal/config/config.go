package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST        string
	DB_PORT        string
	DB_USER        string
	DB_PASSWORD    string
	DB_DATABASE    string
	GRPC           GRPCConfig
	MigrationsPath string
	TokenTTL       time.Duration
}

type GRPCConfig struct {
	Port    int
	Timeout time.Duration
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "8000"
	}
	numPort, err := strconv.Atoi(grpcPort)
	if err != nil {
		panic(err)
	}

	return &Config{
		DB_HOST:     os.Getenv("POSTGRES_HOST"),
		DB_PORT:     os.Getenv("POSTGRES_PORT"),
		DB_USER:     os.Getenv("POSTGRES_USER"),
		DB_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		DB_DATABASE: os.Getenv("POSTGRES_DB"),
		GRPC: GRPCConfig{
			Port:    numPort,
			Timeout: 5 * time.Second,
		},
		MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
		TokenTTL:       24 * time.Hour,
	}
}
