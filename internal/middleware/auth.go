package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"chase_deal/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Authorization header format must be 'Bearer {token}'")
			return
		}

		tokenString := parts[1]

		token, err := utils.ValidateJWT(tokenString)
		if err != nil {
			errorMsg := "Invalid token"

			// Updated error handling for jwt v5
			if err != nil {
				errorMsg := "Invalid token"

				if errors.Is(err, jwt.ErrTokenExpired) {
					errorMsg = "Token expired"
				}

				utils.ErrorResponse(w, http.StatusUnauthorized, errorMsg)
				return
			}

			utils.ErrorResponse(w, http.StatusUnauthorized, errorMsg)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid user in token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
