package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"

type graphQLRequest struct {
	OperationName string `json:"operationName"`
	Query         string `json:"query"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Read body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		// Restore body for next handler
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var gqlReq graphQLRequest
		_ = json.Unmarshal(bodyBytes, &gqlReq)

		// 🔥 Skip auth for public operations
		if SkipAuthMiddleware(gqlReq) {
			next.ServeHTTP(w, r)
			return
		}

		// 🔐 Auth required below
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// skip auth for registration
func SkipAuthMiddleware(req graphQLRequest) bool {
	op := strings.ToLower(req.OperationName)

	// If operationName is provided
	if op == "register" || op == "login" {
		return true
	}

	q := strings.ToLower(req.Query)
	if strings.Contains(q, "mutation") &&
		(strings.Contains(q, "register") || strings.Contains(q, "login")) {
		return true
	}

	return false
}
