package bootstrap

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           int
	PQUser         string
	PQPassword     string
	PQHost         string
	PQPort         int
	PQDatabase     string
	MigrationsPath string
	TokenTTL       int
	RequestTimeout int
	JWTSecret      string
	NatsURL        string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	p := os.Getenv("PORT")
	port := 8080
	if p != "" {
		port, err = strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
	}

	return &Config{
		Port:           port,
		PQUser:         os.Getenv("POSTGRES_USER"),
		PQPassword:     os.Getenv("POSTGRES_PASSWORD"),
		PQHost:         os.Getenv("POSTGRES_HOST"),
		PQPort:         5432,
		PQDatabase:     os.Getenv("POSTGRES_DB"),
		MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
		TokenTTL:       3600 * 24 * 30,
		RequestTimeout: 10,
		JWTSecret:      os.Getenv("JWT_SECRET"),
		NatsURL:        os.Getenv("NATS_URL"),
	}
}
