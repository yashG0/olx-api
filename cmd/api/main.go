package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})
	log.Println("Server is running at PORT 8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
}
