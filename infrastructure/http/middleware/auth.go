package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/fajar-andriansyah/loan-engine/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type contextKey string

const UserContextKey contextKey = "user"

// JWTAuthMiddleware validates JWT token and adds user info to context
func JWTAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sendUnauthorizedResponse(w, "Missing authorization header")
				return
			}

			// Check if header starts with "Bearer "
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				sendUnauthorizedResponse(w, "Invalid authorization header format")
				return
			}

			tokenString := tokenParts[1]

			// Parse and validate token
			claims := &models.JWTClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(viper.GetString("jwt.secret")), nil
			})

			if err != nil {
				log.Error().Err(err).Msg("Failed to parse JWT token")
				sendUnauthorizedResponse(w, "Invalid token")
				return
			}

			if !token.Valid {
				sendUnauthorizedResponse(w, "Invalid token")
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole middleware that checks if user has specific role (for employees)
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserContextKey).(*models.JWTClaims)
			if !ok {
				sendForbiddenResponse(w, "User context not found")
				return
			}

			// Check if user is employee (only employees have roles)
			if user.UserType != "employee" {
				sendForbiddenResponse(w, "Access denied: employee role required")
				return
			}

			// Check if user has required role
			hasRole := false
			for _, role := range roles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				sendForbiddenResponse(w, fmt.Sprintf("Access denied: required role %v", roles))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireUserType middleware that checks if user has specific user type
func RequireUserType(userTypes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserContextKey).(*models.JWTClaims)
			if !ok {
				sendForbiddenResponse(w, "User context not found")
				return
			}

			// Check if user has required user type
			hasUserType := false
			for _, userType := range userTypes {
				if user.UserType == userType {
					hasUserType = true
					break
				}
			}

			if !hasUserType {
				sendForbiddenResponse(w, fmt.Sprintf("Access denied: required user type %v", userTypes))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext extracts user claims from request context
func GetUserFromContext(ctx context.Context) (*models.JWTClaims, error) {
	user, ok := ctx.Value(UserContextKey).(*models.JWTClaims)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}

func sendUnauthorizedResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	response := map[string]interface{}{
		"success":    false,
		"message":    message,
		"error_code": "UNAUTHORIZED",
	}

	json.NewEncoder(w).Encode(models.Response[interface{}]{Data: response})
}

func sendForbiddenResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)

	response := map[string]interface{}{
		"success":    false,
		"message":    message,
		"error_code": "FORBIDDEN",
	}

	json.NewEncoder(w).Encode(models.Response[interface{}]{Data: response})
}
