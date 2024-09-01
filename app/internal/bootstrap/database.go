package bootstrap

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDatabaseConnection(cfg *Config) *sql.DB {
	db, err := sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.PQUser, cfg.PQPassword, cfg.PQHost, cfg.PQPort, cfg.PQDatabase))
	if err != nil {
		panic(err)
	}
	return db
}

func MigrateDatabase(db *sql.DB, cfg *Config) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.PQUser, cfg.PQPassword, cfg.PQHost, cfg.PQPort, cfg.PQDatabase)

	m, err := migrate.New("file://"+cfg.MigrationsPath, dsn)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}
