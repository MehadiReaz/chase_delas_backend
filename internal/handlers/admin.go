package handlers

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"net/http"

	"chase_deal/internal/models"
	"chase_deal/pkg/utils"
)

type AdminHandler struct {
	DB *sql.DB
}

// GetAllUsers - Superadmin/Admin can list all users
func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, email, role, created_at FROM users`
	rows, err := h.DB.Query(query)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt)
		if err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to read user data")
			return
		}
		users = append(users, user)
	}

	utils.SuccessResponse(w, http.StatusOK, "Users retrieved", users)
}

// UpdateUserRole - Superadmin can update roles
func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	var updateReq struct {
		UserID  string `json:"user_id"`
		NewRole string `json:"new_role"`
	}

	if err := utils.ParseJSON(r, &updateReq); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Validate role
	validRoles := map[string]bool{"admin": true, "vendor": true, "user": true}
	if !validRoles[updateReq.NewRole] {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid role specified")
		return
	}

	_, err := h.DB.Exec(
		"UPDATE users SET role = ? WHERE id = ?",
		updateReq.NewRole,
		updateReq.UserID,
	)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update role")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Role updated successfully", nil)
}

// DeactivateUser - Admin can deactivate users
func (h *AdminHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	_, err := h.DB.Exec("UPDATE users SET active = false WHERE id = ?", userID)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to deactivate user")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "User deactivated", nil)
}

func (h *AdminHandler) CreateShop(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		OwnerID  string `json:"owner_id"`
		Priority int    `json:"priority"`
	}

	if err := utils.ParseJSON(r, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Verify owner exists and is shop owner role
	var role string
	err := h.DB.QueryRow("SELECT role FROM users WHERE id = ?", req.OwnerID).Scan(&role)
	if err != nil || role != "vendor" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid shop owner")
		return
	}

	shopID := xid.New().String()
	_, err = h.DB.Exec(
		`INSERT INTO shops (id, name, owner_id, priority) 
        VALUES (?, ?, ?, ?)`,
		shopID, req.Name, req.OwnerID, req.Priority,
	)

	utils.SuccessResponse(w, http.StatusCreated, "Shop created", map[string]string{"id": shopID})
}
