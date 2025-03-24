package handlers

import (
	"database/sql"
	"net/http"

	"chase_deal/internal/models"
	"chase_deal/pkg/utils"
	"github.com/gorilla/mux"
)

type ShopHandler struct {
	DB *sql.DB
}

func (h *ShopHandler) UpdateShop(w http.ResponseWriter, r *http.Request) {
	shopID := mux.Vars(r)["id"]
	var req models.ShopUpdateRequest

	if err := utils.ParseJSON(r, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Verify ownership
	var ownerID string
	err := h.DB.QueryRow("SELECT owner_id FROM shops WHERE id = ?", shopID).
		Scan(&ownerID)
	if err != nil || ownerID != utils.GetUserID(r) {
		utils.ErrorResponse(w, http.StatusForbidden, "Not the shop owner")
		return
	}

	_, err = h.DB.Exec(
		`UPDATE shops SET 
        name = ?, latitude = ?, longitude = ? 
        WHERE id = ?`,
		req.Name, req.Latitude, req.Longitude, shopID,
	)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update shop")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Shop updated", nil)
}

func (h *ShopHandler) GetShopDetails(w http.ResponseWriter, r *http.Request) {
	shopID := mux.Vars(r)["id"]

	var shop models.Shop
	err := h.DB.QueryRow(
		`SELECT id, name, owner_id, priority, latitude, longitude, 
        is_active, created_at FROM shops WHERE id = ?`,
		shopID,
	).Scan(
		&shop.ID, &shop.Name, &shop.OwnerID, &shop.Priority,
		&shop.Latitude, &shop.Longitude, &shop.IsActive, &shop.CreatedAt,
	)

	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Shop not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Shop details", shop)
}
