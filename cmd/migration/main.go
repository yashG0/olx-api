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
	// if argument less then 2 then no usage of migration main file, exit!
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down>")
	}
	// must load config data from envirnmental variable otherwise exit the program
	cfg := config.MustLoad()

	// make new migration instance using migration reading file and database url
	m, err := migrate.New(
		"file://migrations",
		cfg.DBUrl)

	// if migration failed to initialize exit with error
	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}

	// handle the migration logic from argument if invalid argument make error
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
