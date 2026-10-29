package httpx

import (
	"encoding/json"
	"net/http"
)

type CodeType string

const (
	CodeInvalidID        CodeType = "invalid id"
	CodeNotFound         CodeType = "not found"
	CodeInternalError    CodeType = "internal error"
	CodeMalformedJson    CodeType = "malformed json"
	CodeValidationFailed CodeType = "validation failed"
	CodeUnauthenticate   CodeType = "unauthenticated"
	CodeForbidden        CodeType = "forbidden"
	CodeConflict         CodeType = "conflict"
	CodeRateLimited      CodeType = "rate limited"
)

type errorPayload struct {
	Code    CodeType `json:"code"`
	Message string   `json:"message"`
	Field   string   `json:"field,omitempty"`
}

type errorEnvelope struct {
	Error errorPayload `json:"error"`
}

func Error(w http.ResponseWriter, message string, code CodeType, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorEnvelope{Error: errorPayload{Code: code, Message: message}})
}

func ValidationError(w http.ResponseWriter, message string, code CodeType, status int, field string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorEnvelope{Error: errorPayload{Code: code, Message: message, Field: field}})
}
