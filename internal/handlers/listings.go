package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ListingHandler struct {
	db *sql.DB
}

func NewListingHandler(db *sql.DB) *ListingHandler {
	return &ListingHandler{db: db}
}

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (lh ListingHandler) Get(w http.ResponseWriter, r *http.Request) {
	const (
		defaultLimit = 20
		maxLimit     = 100
	)

	limit := defaultLimit
	if val := r.URL.Query().Get("limit"); val != "" {
		n, err := strconv.Atoi(val)
		if err != nil || n < 1 {
			http.Error(
				w,
				"limit must be a positive integer",
				http.StatusBadRequest,
			)
			return
		}

		limit = min(n, maxLimit)
	}

	rows, err := lh.db.Query(`
		SELECT id, title, description, price, city, created_at, updated_at
		FROM listings
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		log.Printf("query.listings: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	listings := make([]listing, 0, limit)

	for rows.Next() {
		var l listing
		if err := rows.Scan(
			&l.ID,
			&l.Title,
			&l.Description,
			&l.Price,
			&l.City,
			&l.CreatedAt,
			&l.UpdatedAt,
		); err != nil {
			log.Printf("rows.scan: %v", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows.err: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(listings)
}

func (lh ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, err := lh.db.Exec("DELETE FROM listings WHERE id=$1", id)
	if err != nil {
		log.Printf("delete: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
