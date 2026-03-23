package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"

// AuthMiddleware validates the Bearer token when present and injects the user ID
// into the request context. Requests without a token pass through unauthenticated;
// resolvers are responsible for rejecting requests that require auth.
func AuthMiddleware(ts *TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
				if claims, err := ts.ValidateToken(tokenStr); err == nil {
					ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
					r = r.WithContext(ctx)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
