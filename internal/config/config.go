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
		grpcPort = "50051"
	}
	numPort, err := strconv.Atoi(grpcPort)
	if err != nil {
		panic(err)
	}

	return &Config{
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_DATABASE: os.Getenv("DB_DATABASE"),
		GRPC: GRPCConfig{
			Port:    numPort,
			Timeout: 5 * time.Second,
		},
		MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
		TokenTTL:       24 * time.Hour,
	}
}
