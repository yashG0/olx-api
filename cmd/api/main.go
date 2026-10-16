package main

import (
	"log"
	"net/http"
	"time"

	"github.com/yashG0/olx-api/internal/config"
)

func main() {
	cfg := config.MustLoad()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	serverConfig := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Println("Server is running at PORT", cfg.Port)
	err := serverConfig.ListenAndServe()
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
}
