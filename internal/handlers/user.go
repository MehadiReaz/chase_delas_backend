package handlers

import (
	"database/sql"
	"net/http"

	"chase_deal/internal/models"
	"chase_deal/pkg/utils"
)

type UserHandler struct {
	DB *sql.DB
}

// GetProfile - User gets their profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r)

	var user models.User
	err := h.DB.QueryRow(
		"SELECT id, email, role, created_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt)

	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Profile retrieved", user)
}

// UpdateProfile - User updates their profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r)
	var updateData struct {
		Email string `json:"email"`
	}

	if err := utils.ParseJSON(r, &updateData); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	_, err := h.DB.Exec(
		"UPDATE users SET email = ? WHERE id = ?",
		updateData.Email,
		userID,
	)

	if err != nil {
		utils.ErrorResponse(w, http.StatusConflict, "Email already in use")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Profile updated", nil)
}
