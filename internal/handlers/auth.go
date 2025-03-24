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

	// Get requester role only if Authorization header exists
	var requesterRole string
	if r.Header.Get("Authorization") != "" {
		requesterRole = utils.GetUserRole(r)
	}

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	// Authorization logic
	if user.Role == "admin" && requesterRole != "superadmin" {
		utils.ErrorResponse(w, http.StatusForbidden, "Only superadmins can create admins")
		return
	}

	user.ID = xid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Could not hash password")
		return
	}

	_, err = h.DB.Exec(
		"INSERT INTO users (id, email, password, role) VALUES (?, ?, ?, ?)",
		user.ID, user.Email, string(hashedPassword), user.Role,
	)

	if err != nil {
		if utils.IsDuplicateError(err) {
			utils.ErrorResponse(w, http.StatusConflict, "User already exists")
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Registration failed: "+err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "User created successfully", map[string]string{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (h *AuthHandler) UpdateCredentials(w http.ResponseWriter, r *http.Request) {
	var updateReq models.UpdateCredentialsRequest
	if err := utils.ParseJSON(r, &updateReq); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := utils.GetUserID(r)
	requesterRole := utils.GetUserRole(r)

	// Get existing user
	var existingUser models.User
	err := h.DB.QueryRow(
		"SELECT id, email, password, role FROM users WHERE id = ?",
		userID,
	).Scan(&existingUser.ID, &existingUser.Email, &existingUser.Password, &existingUser.Role)

	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	// Superadmin can update anyone, others can only update themselves
	if requesterRole != "superadmin" && existingUser.ID != userID {
		utils.ErrorResponse(w, http.StatusForbidden, "You can only update your own credentials")
		return
	}

	// Update email
	if updateReq.Email != "" {
		if updateReq.Email == existingUser.Email {
			utils.ErrorResponse(w, http.StatusBadRequest, "New email same as current email")
			return
		}

		_, err = h.DB.Exec(
			"UPDATE users SET email = ? WHERE id = ?",
			updateReq.Email, userID,
		)
		if err != nil {
			utils.ErrorResponse(w, http.StatusConflict, "Email already in use")
			return
		}
	}

	// Update password
	if updateReq.NewPassword != "" {
		if requesterRole != "superadmin" {
			// Verify old password for non-superadmins
			err = bcrypt.CompareHashAndPassword(
				[]byte(existingUser.Password),
				[]byte(updateReq.OldPassword),
			)
			if err != nil {
				utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid current password")
				return
			}
		}

		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateReq.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Could not hash password")
			return
		}

		_, err = h.DB.Exec(
			"UPDATE users SET password = ? WHERE id = ?",
			string(newHashedPassword), userID,
		)
		if err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Password update failed")
			return
		}
	}

	utils.SuccessResponse(w, http.StatusOK, "Credentials updated successfully", nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.LoginRequest
	if err := utils.ParseJSON(r, &creds); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var storedUser models.User
	// Get user with role information
	err := h.DB.QueryRow(
		"SELECT id, email, password, role FROM users WHERE email = ?",
		creds.Email,
	).Scan(&storedUser.ID, &storedUser.Email, &storedUser.Password, &storedUser.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(storedUser.Password),
		[]byte(creds.Password),
	); err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT with role claim
	token, err := utils.GenerateJWT(storedUser.ID, storedUser.Role)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Could not generate token: "+err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Login successful", map[string]string{
		"token": token,
		"id":    storedUser.ID,
		"email": storedUser.Email,
		"role":  storedUser.Role,
	})
}
