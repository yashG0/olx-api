package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey int

const (
	requestId    string = "X-Request-ID"
	requestCtxId ctxKey = iota
)

// added the middleware for listing handler where we receive or generate the request id for user and forward to listing handler
func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestId)
		if id == "" {
			id = uuid.NewString()
		}

		w.Header().Add(requestId, id)
		ctx := context.WithValue(r.Context(), requestCtxId, id)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func RequestIDFromContext(ctx context.Context) string {
	reqId := ctx.Value(requestCtxId).(string)
	return reqId
}
