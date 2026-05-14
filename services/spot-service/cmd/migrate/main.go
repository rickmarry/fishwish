package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env.local")

	action := flag.String("action", "up", "migration action: up, down, version")
	flag.Parse()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	m, err := migrate.New(
		"file://internal/repository/migrations",
		"postgres://"+dbURL[11:],
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	switch *action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration up failed: %v", err)
		}
		log.Println("migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration down failed: %v", err)
		}
		log.Println("migrations rolled back successfully")
	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			log.Printf("no migrations applied yet")
			return
		}
		fmt.Printf("version: %d (dirty: %t)\n", v, dirty)
	default:
		log.Fatalf("unknown action: %s", *action)
	}
}
