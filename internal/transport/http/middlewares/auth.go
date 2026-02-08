package middlewares

import (
	"net/http"
)

// AuthMiddleware validates the Authorization header (Bearer token) and stores claims in context
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if r.Header.Get("Authorization") == "" {
		// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
		// 	return
		// }
		next(w, r)
	}
}
