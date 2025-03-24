package handlers

import (
	"net/http"

	"chase_deal/internal/models"
	"chase_deal/pkg/utils"
)

//type OfferHandler struct {
//	DB *sql.DB
//}

func (h *OfferHandler) GetOffers(w http.ResponseWriter, r *http.Request) {
	query := `SELECT o.id, o.shop_id, o.title, o.description, 
             o.discount_value, o.start_date, o.end_date, 
             o.image_url, o.is_featured, o.created_at 
             FROM offers o WHERE o.end_date > CURRENT_TIMESTAMP`

	rows, err := h.DB.Query(query)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to get offers")
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

	utils.SuccessResponse(w, http.StatusOK, "Offers retrieved", offers)
}
