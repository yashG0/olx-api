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
	ctx := r.Context()
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

	rows, err := lh.db.QueryContext(ctx, `
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
	ctx := r.Context()
	id := r.PathValue("id")
	_, err := lh.db.ExecContext(ctx, "DELETE FROM listings WHERE id=$1", id)
	if err != nil {
		log.Printf("delete: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (lh ListingHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	res := lh.db.QueryRowContext(ctx, "SELECT id, title, description, price, city, created_at, updated_at FROM listings WHERE id = $1", id)
	var data listing
	err := res.Scan(
		&data.ID,
		&data.Title,
		&data.Description,
		&data.Price,
		&data.City,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		log.Printf("listing.getid: %v", err)
		http.Error(w, "Failed to parser data!", http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (lh ListingHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	var data listing
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("listing.update.decode: %v", err)
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	res, err := lh.db.ExecContext(ctx, `UPDATE listings SET title = $1, description = $2, price = $3, city = $4, updated_at = NOW()
		WHERE id = $5
	`, data.Title, data.Description, data.Price, data.City, id)

	if err != nil {
		log.Printf("listing.update: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("listing.update.rows_affected: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, "listing not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
