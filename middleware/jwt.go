package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/abdulkarimogaji/invoGenius/services/token"
)

type apiResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type contextKey string

const UserIDKey contextKey = "user_id"

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := extractToken(r)
		if err != nil {
			errorResponse(err, w)
			return
		}
		claims, err := token.ValidateToken(tokenString)
		if err != nil {
			errorResponse(err, w)
			return
		}
		userID, err := claims.GetSubject()
		if err != nil {
			errorResponse(err, w)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	// Authorization header should be in the format "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}

func errorResponse(err error, w http.ResponseWriter) {
	response := apiResponse{
		Error:   true,
		Message: err.Error(),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(jsonResponse)
}
