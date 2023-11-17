package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// TrackRequestMiddleware adds an unique id to each request for it to be tracked
// afterwards in logs if needed.
func TrackRequestMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := r.Context()

	requestID := uuid.New().String()
	ctx = context.WithValue(ctx, "requestID", requestID)

	r = r.WithContext(ctx)

	next(w, r)
}
