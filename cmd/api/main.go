package main

import (
	"log"
	"net/http"
	"time"

	"github.com/yashG0/olx-api/internal/config"
	"github.com/yashG0/olx-api/internal/db"
	"github.com/yashG0/olx-api/internal/handlers"
)

func main() {
	cfg := config.MustLoad()
	_, err := db.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}
	mux := http.NewServeMux()

	// adding routes
	mux.HandleFunc("GET /health", handlers.Health)

	// creating configuration setup
	serverConfig := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Println("Server is running at PORT", cfg.Port)
	err = serverConfig.ListenAndServe()
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
}
