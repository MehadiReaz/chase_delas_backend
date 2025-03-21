package handlers

import (
	"database/sql"
	"net/http"

	"chase_deal/internal/models"
	"chase_deal/pkg/utils"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generate unique ID
	user.ID = xid.New().String()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Could not hash password")
		return
	}

	// Store user in database
	_, err = h.DB.Exec("INSERT INTO users (id, email, password) VALUES (?, ?, ?)",
		user.ID, user.Email, string(hashedPassword))

	if err != nil {
		// Handle duplicate email error specifically
		if utils.IsDuplicateError(err) {
			utils.ErrorResponse(w, http.StatusConflict, "User already exists")
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Registration failed")
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "User created successfully", nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.LoginRequest
	if err := utils.ParseJSON(r, &creds); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var storedUser models.User
	err := h.DB.QueryRow("SELECT id, email, password FROM users WHERE email = ?", creds.Email).
		Scan(&storedUser.ID, &storedUser.Email, &storedUser.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password)); err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := utils.GenerateJWT(storedUser.ID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Login successful", map[string]string{
		"token": token,
	})
}
