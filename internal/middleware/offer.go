package middleware

import (
	"database/sql"
	"net/http"

	"chase_deal/pkg/utils"
	"github.com/gorilla/mux"
)

type OfferOwnerMiddleware struct {
	DB *sql.DB
}

func NewOfferOwnerMiddleware(db *sql.DB) *OfferOwnerMiddleware {
	return &OfferOwnerMiddleware{DB: db}
}

func (m *OfferOwnerMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offerID := mux.Vars(r)["id"]
		userID := utils.GetUserID(r)

		var shopID string
		err := m.DB.QueryRow("SELECT shop_id FROM offers WHERE id = ?", offerID).
			Scan(&shopID)
		if err != nil {
			utils.ErrorResponse(w, http.StatusNotFound, "Offer not found")
			return
		}

		var ownerID string
		err = m.DB.QueryRow("SELECT owner_id FROM shops WHERE id = ?", shopID).
			Scan(&ownerID)
		if err != nil || ownerID != userID {
			utils.ErrorResponse(w, http.StatusForbidden, "Not authorized for this offer")
			return
		}

		next.ServeHTTP(w, r)
	})
}
