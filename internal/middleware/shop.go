package middleware

import (
	"database/sql"
	"net/http"

	"chase_deal/pkg/utils"
	"github.com/gorilla/mux"
)

type ShopOwnerMiddleware struct {
	DB *sql.DB
}

func NewShopOwnerMiddleware(db *sql.DB) *ShopOwnerMiddleware {
	return &ShopOwnerMiddleware{DB: db}
}

func (m *ShopOwnerMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shopID := mux.Vars(r)["shopId"]
		userID := utils.GetUserID(r)

		var ownerID string
		err := m.DB.QueryRow(
			"SELECT owner_id FROM shops WHERE id = ?",
			shopID,
		).Scan(&ownerID)

		if err != nil {
			if err == sql.ErrNoRows {
				utils.ErrorResponse(w, http.StatusNotFound, "Shop not found")
			} else {
				utils.ErrorResponse(w, http.StatusInternalServerError, "Database error")
			}
			return
		}

		if ownerID != userID {
			utils.ErrorResponse(w, http.StatusForbidden, "Not the shop owner")
			return
		}

		next.ServeHTTP(w, r)
	})
}
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := utils.GetUserRole(r)
			for _, role := range roles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			utils.ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		})
	}
}
