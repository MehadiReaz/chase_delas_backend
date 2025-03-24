package utils

import (
	"chase_deal/internal/models"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

// JWT functions
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func GetUserRole(r *http.Request) string {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return ""
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || token == nil || !token.Valid {
		return ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if role, ok := claims["role"].(string); ok {
			return role
		}
	}
	return ""
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
}

// JSON parsing
func ParseJSON(r *http.Request, v interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("invalid content type, expected application/json")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("error decoding request body: %v", err)
	}

	return nil
}

// Response helpers
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(models.ErrorResponseModel{
		Status:  statusCode,
		Message: message,
	})
	if err != nil {
		return
	}
}

func SuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(models.SuccessResponseModel{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
	if err != nil {
		return
	}
}

func IsDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

func GetUserID(r *http.Request) string {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return ""
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || token == nil || !token.Valid {
		return ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userID, ok := claims["user_id"].(string); ok {
			return userID
		}
	}
	return ""
}
