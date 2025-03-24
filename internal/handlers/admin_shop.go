package handlers

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"

	"chase_deal/internal/models"
	"chase_deal/pkg/utils"
	"github.com/rs/xid"
)

type AdminShopHandler struct {
	DB *sql.DB
}

func (h *AdminShopHandler) CreateShop(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      string  `json:"name"`
		OwnerID   string  `json:"owner_id"`
		Priority  int     `json:"priority"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
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
		`INSERT INTO shops 
        (id, name, owner_id, priority, latitude, longitude) 
        VALUES (?, ?, ?, ?, ?, ?)`,
		shopID, req.Name, req.OwnerID, req.Priority,
		req.Latitude, req.Longitude,
	)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create shop")
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Shop created",
		map[string]string{"id": shopID})
}

func (h *AdminShopHandler) UpdateShopPriority(w http.ResponseWriter, r *http.Request) {
	shopID := mux.Vars(r)["id"]
	var req struct {
		Priority int `json:"priority"`
	}

	if err := utils.ParseJSON(r, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	_, err := h.DB.Exec(
		"UPDATE shops SET priority = ? WHERE id = ?",
		req.Priority, shopID,
	)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update priority")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Priority updated", nil)
}

func (h *AdminShopHandler) ListShops(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name, owner_id, priority, latitude, longitude, 
             is_active, created_at FROM shops`

	rows, err := h.DB.Query(query)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch shops")
		return
	}
	defer rows.Close()

	var shops []models.Shop
	for rows.Next() {
		var shop models.Shop
		err := rows.Scan(
			&shop.ID, &shop.Name, &shop.OwnerID, &shop.Priority,
			&shop.Latitude, &shop.Longitude, &shop.IsActive, &shop.CreatedAt,
		)
		if err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to read shop data")
			return
		}
		shops = append(shops, shop)
	}

	utils.SuccessResponse(w, http.StatusOK, "Shops retrieved", shops)
}
