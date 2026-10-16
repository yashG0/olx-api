package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	mux := http.NewServeMux()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	port := os.Getenv("PORT")
	serverConfig := http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Println("Server is running at PORT", port)
	err = serverConfig.ListenAndServe()
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
}
