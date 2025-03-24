package middleware

import (
	"net/http"

	"chase_deal/pkg/utils"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := utils.GetUserRole(r)
		if role != "admin" {
			utils.ErrorResponse(w, http.StatusForbidden, "Admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
