package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/yashG0/olx-api/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down>")
	}
	cfg := config.MustLoad()

	m, err := migrate.New(
		"file://migrations",
		cfg.DBUrl)

	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}

	switch os.Args[1] {
	case "up":
		err := m.Steps(1)
		if err != nil {
			log.Fatalf("migration.up: %v", err)
		}
	case "down":
		err := m.Steps(-1)
		if err != nil {
			log.Fatalf("migration.down: %v", err)
		}
	default:
		log.Fatalf("Unknown command %s", os.Args[1])
	}
}
