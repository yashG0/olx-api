package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/yashG0/olx-api/internal/config"
	"github.com/yashG0/olx-api/internal/db"
	"github.com/yashG0/olx-api/internal/handlers"
)

func main() {
	// loading main configuration or direct program exit
	cfg := config.MustLoad()

	// we handle database connection if connect ok else we get error
	db, err := db.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

	// Configure structured logger here
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	// Register Handler into New slog instance
	logger := slog.New(logHandler)

	// we have created new Router/Mux
	mux := http.NewServeMux()

	// adding new Route in new mux
	lh := handlers.NewListingHandler(db, logger)
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /listing", lh.Get)
	mux.HandleFunc("GET /listing/{id}", lh.GetById)
	mux.HandleFunc("DELETE /listing/{id}", lh.Delete)

	// creating configuration setup
	serverConfig := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Println("Server is running at PORT", cfg.Port)

	// Listeing our Router from serverconfig
	err = serverConfig.ListenAndServe()
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
}
