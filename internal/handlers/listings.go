package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/yashG0/olx-api/internal/httpx"
	"github.com/yashG0/olx-api/internal/middleware"
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
	// create a request context for User Request
	ctx := r.Context()

	// creating constant variable for limit query parameter for pagination of listings
	const (
		defaultLimit = 20
		maxLimit     = 100
	)

	// configure the limit parameter gracefully handling it
	limit := defaultLimit
	if val := r.URL.Query().Get("limit"); val != "" {
		n, err := strconv.Atoi(val)
		if err != nil || n < 1 {
			httpx.Error(
				w,
				"limit must be a positive integer",
				httpx.CodeValidationFailed,
				http.StatusBadRequest,
			)
			return
		}

		limit = min(n, maxLimit)
	}

	// making a DB query for getting the listing elements with passing the context
	rows, err := lh.db.QueryContext(ctx, `
		SELECT id, title, description, price, city, created_at, updated_at
		FROM listings
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		lh.logger.Info("Failed to fetched Result!", "Error", err)
		httpx.Error(w, "database error", httpx.CodeInternalError, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// reading the result rows in in-memory variable
	listings := make([]listing, 0, limit)

	// we can go forward using Next() function to move on
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

	// checking the error may arise while scanning the each rows
	err = rows.Err()
	if err != nil {
		lh.logger.Error("listing.get.parse: Failed to scan rows!", "Error", err)
		httpx.Error(w, "Failed to parse row!", httpx.CodeInternalError, http.StatusInternalServerError)
		return
	}

	lh.logger.Info("listing fetched", "total", len(listings))

	// if everything is successfull we will make a the Response to user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(listings)
}

func (lh ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// again, create a request context of user we will pass while DB query
	ctx := r.Context()
	requestID := middleware.RequestIDFromContext(ctx)

	// getting the path parameter
	id := r.PathValue("id")

	// query for the deleting the listing by id
	_, err := lh.db.ExecContext(ctx, "DELETE FROM listings WHERE id=$1", id)

	// if failed to delete the listing arise the error
	if err != nil {
		lh.logger.Error("failed to delete record in listing!", "listing_id", id, "error", err, "request-id", requestID)
		httpx.Error(w, "failed to delete record in listing!", httpx.CodeInternalError, http.StatusInternalServerError)
		return
	}

	// make a header to Response to
	w.WriteHeader(http.StatusNoContent)
}

func (lh ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqId := middleware.RequestIDFromContext(ctx)
	defer r.Body.Close()

	var l CreatelistingRequest
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		lh.logger.Error("Failed to decode!", "request-id", reqId)
		httpx.Error(w, err.Error(), httpx.CodeMalformedJson, http.StatusBadRequest)
		return
	}

	row := lh.db.QueryRowContext(
		ctx,
		"INSERT INTO listings (title,description,price,city) VALUES ($1,$2,$3,$4) RETURNING id,title,description,created_at",
		l.Title, l.Description, l.Price, l.City,
	)
	var out CreateListingResponse

	err = row.Scan(&out.ID, &out.Title, &out.Description, &out.CreatedAt)
	if err != nil {
		lh.logger.Error("listing.create.scan: Failed to insert!", "request_id", reqId)
		httpx.Error(w, err.Error(), httpx.CodeInternalError, http.StatusInternalServerError)
		return
	}
	lh.logger.Info("listing.created.success: New listing Created", "request-id", reqId, "listing-id", out.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(out)
}
