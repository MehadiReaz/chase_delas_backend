package middleware

import (
	"net/http"

	"chase_deal/pkg/utils"
)

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := utils.GetUserRole(r)

			for _, role := range allowedRoles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			utils.ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		})
	}
}
