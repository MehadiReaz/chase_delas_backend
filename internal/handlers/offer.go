package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"chase_deal/internal/models"
	"chase_deal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
)

type OfferHandler struct {
	DB *sql.DB
}

func (h *OfferHandler) CreateOffer(w http.ResponseWriter, r *http.Request) {
	var offer models.Offer
	if err := utils.ParseJSON(r, &offer); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid offer data")
		return
	}

	// Verify shop ownership
	var ownerID string
	err := h.DB.QueryRow("SELECT owner_id FROM shops WHERE id = ?", offer.ShopID).
		Scan(&ownerID)
	if err != nil || ownerID != utils.GetUserID(r) {
		utils.ErrorResponse(w, http.StatusForbidden, "Not authorized for this shop")
		return
	}

	offer.ID = xid.New().String()
	offer.CreatedAt = time.Now().UTC()

	_, err = h.DB.Exec(
		`INSERT INTO offers 
        (id, shop_id, title, description, discount_value, 
         start_date, end_date, image_url, is_featured, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		offer.ID, offer.ShopID, offer.Title, offer.Description,
		offer.DiscountValue, offer.StartDate, offer.EndDate,
		offer.ImageURL, offer.IsFeatured, offer.CreatedAt,
	)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create offer")
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Offer created", offer)
}

func (h *OfferHandler) UpdateOffer(w http.ResponseWriter, r *http.Request) {
	offerID := mux.Vars(r)["id"]
	var updateData models.OfferUpdateRequest

	if err := utils.ParseJSON(r, &updateData); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid update data")
		return
	}

	// Verify offer ownership through shop
	var shopID string
	err := h.DB.QueryRow("SELECT shop_id FROM offers WHERE id = ?", offerID).Scan(&shopID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Offer not found")
		return
	}

	var ownerID string
	err = h.DB.QueryRow("SELECT owner_id FROM shops WHERE id = ?", shopID).Scan(&ownerID)
	if err != nil || ownerID != utils.GetUserID(r) {
		utils.ErrorResponse(w, http.StatusForbidden, "Not authorized to update this offer")
		return
	}

	_, err = h.DB.Exec(
		`UPDATE offers SET
        title = ?, description = ?, discount_value = ?,
        start_date = ?, end_date = ?
        WHERE id = ?`,
		updateData.Title, updateData.Description, updateData.DiscountValue,
		updateData.StartDate, updateData.EndDate, offerID,
	)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update offer")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Offer updated", nil)
}

func (h *OfferHandler) GetOffer(w http.ResponseWriter, r *http.Request) {
	offerID := mux.Vars(r)["id"]

	var offer models.Offer
	err := h.DB.QueryRow(
		`SELECT id, shop_id, title, description, discount_value,
        start_date, end_date, image_url, is_featured, created_at
        FROM offers WHERE id = ?`,
		offerID,
	).Scan(
		&offer.ID, &offer.ShopID, &offer.Title, &offer.Description,
		&offer.DiscountValue, &offer.StartDate, &offer.EndDate,
		&offer.ImageURL, &offer.IsFeatured, &offer.CreatedAt,
	)

	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Offer not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Offer retrieved", offer)
}

func (h *OfferHandler) ListOffers(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, shop_id, title, description, discount_value,
             start_date, end_date, image_url, is_featured, created_at
             FROM offers WHERE end_date > ?`

	rows, err := h.DB.Query(query, time.Now().UTC())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch offers")
		return
	}
	defer rows.Close()

	var offers []models.Offer
	for rows.Next() {
		var offer models.Offer
		err := rows.Scan(
			&offer.ID, &offer.ShopID, &offer.Title, &offer.Description,
			&offer.DiscountValue, &offer.StartDate, &offer.EndDate,
			&offer.ImageURL, &offer.IsFeatured, &offer.CreatedAt,
		)
		if err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Error reading offer data")
			return
		}
		offers = append(offers, offer)
	}

	utils.SuccessResponse(w, http.StatusOK, "Active offers retrieved", offers)
}

func (h *OfferHandler) ToggleFeatured(w http.ResponseWriter, r *http.Request) {
	offerID := mux.Vars(r)["id"]

	// Verify ownership
	var shopID string
	err := h.DB.QueryRow("SELECT shop_id FROM offers WHERE id = ?", offerID).Scan(&shopID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Offer not found")
		return
	}

	var ownerID string
	err = h.DB.QueryRow("SELECT owner_id FROM shops WHERE id = ?", shopID).Scan(&ownerID)
	if err != nil || ownerID != utils.GetUserID(r) {
		utils.ErrorResponse(w, http.StatusForbidden, "Not authorized to modify this offer")
		return
	}

	_, err = h.DB.Exec(
		"UPDATE offers SET is_featured = NOT is_featured WHERE id = ?",
		offerID,
	)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to toggle feature status")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Feature status updated", nil)
}
