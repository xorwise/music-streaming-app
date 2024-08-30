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

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, dbHost, dbPort, dbName)
	fmt.Println(dsn)

	m, err := migrate.New("file://"+migrationsPath, dsn)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

}
