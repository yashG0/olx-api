package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{db: db, logger: logger}
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
		lh.logger.Info("Failed to fetched Result!", "Error", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	listings := make([]listing, 0, limit)

	for rows.Next() {
		var l listing
		err := rows.Scan(
			&l.ID,
			&l.Title,
			&l.Description,
			&l.Price,
			&l.City,
			&l.CreatedAt,
			&l.UpdatedAt,
		)
		if err != nil {
			lh.logger.Error("Failed to parse Result!", "Error", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		listings = append(listings, l)
	}

	err = rows.Err()
	if err != nil {
		lh.logger.Error("Failed to parse each Result", "Error", err)
		http.Error(w, "Internal Error!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(listings)
}

func (lh ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	_, err := lh.db.ExecContext(ctx, "DELETE FROM listing WHERE id=$1", id)
	if err != nil {
		lh.logger.Error("failed to delete record in listing!", "listing_id", id, "error", err)
		http.Error(w, "failed to delete record in listing!", http.StatusInternalServerError)
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
		lh.logger.Error("Failed to parse data!", "Error", err, "Id", id)
		http.Error(w, "Failed to parser data!", http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
