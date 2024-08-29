package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	var (
		migrationsPath  string = "./migrations"
		migrationsTable string = "migrations"
	)

	flag.StringVar(&migrationsPath, "migrations-path", migrationsPath, "Migrations path")
	flag.StringVar(&migrationsTable, "migrations-table", migrationsTable, "Migrations table")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_DATABASE")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, dbHost, dbPort, dbName)

	m, err := migrate.New("file://"+migrationsPath, dsn)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

}
