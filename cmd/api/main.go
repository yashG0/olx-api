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
	"github.com/yashG0/olx-api/internal/middleware"
)

func main() {
	// loading main configuration or direct program exit -
	cfg := config.MustLoad()

	// we handle database connection if connect then ok otherwise we will get error -
	db, err := db.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

	// Open log file.
	logFile, err := os.OpenFile(
		"logger/data.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatalf("main.logger.open: %v", err)
	}
	defer logFile.Close()

	// Configure structured logger (file only).
	logHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})

	logger := slog.New(logHandler)

	// we have created new Router/Mux -
	mux := http.NewServeMux()

	// adding new Route in new mux -
	lh := handlers.NewListingHandler(db, logger)
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /listing", lh.Get)
	mux.HandleFunc("DELETE /listing/{id}", lh.Delete)
	mux.HandleFunc("POST /listing", lh.Create)

	// wiring the middleware with main route -
	handler := middleware.RequestId(mux)

	// creating configuration setup -
	serverConfig := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Println("Server is running at PORT", cfg.Port)

	// Listeing our Router from serverconfig -
	err = serverConfig.ListenAndServe()
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
}
