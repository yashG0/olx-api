package handlers

import (
	"fmt"
	"strings"
	"time"
)

type CreatelistingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	City        string `json:"city"`
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func (req CreatelistingRequest) validateRequestData() error {
	if strings.TrimSpace(req.Title) == "" {
		return &ValidationError{Field: "title", Message: "must not be empty!"}
	}
	if strings.TrimSpace(req.Description) == "" {
		return &ValidationError{Field: "description", Message: "must not be empty!"}
	}
	if strings.TrimSpace(req.City) == "" {
		return &ValidationError{Field: "city", Message: "must not be empty!"}
	}
	if req.Price <= 0 {
		return &ValidationError{Field: "price", Message: "cannot be less then or equal to zero"}
	}
	return nil
}

type CreateListingResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
