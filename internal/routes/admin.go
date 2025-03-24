package routes

import (
	"net/http"

	"chase_deal/internal/handlers"
	"chase_deal/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterAdminRoutes(
	router *mux.Router,
	adminHandler *handlers.AdminHandler,
	adminShopHandler *handlers.AdminShopHandler,
) {
	adminRouter := router.PathPrefix("/admin").Subrouter()

	// Common middleware for all admin routes
	adminRouter.Use(middleware.JWTAuth)
	adminRouter.Use(middleware.RoleMiddleware("superadmin"))

	// User management routes
	adminRouter.HandleFunc("/users", adminHandler.GetAllUsers).Methods(http.MethodGet)
	adminRouter.HandleFunc("/users/{id}/role", adminHandler.UpdateUserRole).Methods(http.MethodPut)
	adminRouter.HandleFunc("/users/{id}", adminHandler.DeactivateUser).Methods(http.MethodDelete)

	// Shop management routes
	adminRouter.HandleFunc("/shops", adminShopHandler.ListShops).Methods(http.MethodGet)
	adminRouter.HandleFunc("/shops", adminShopHandler.CreateShop).Methods(http.MethodPost)
	adminRouter.HandleFunc("/shops/{id}/priority", adminShopHandler.UpdateShopPriority).Methods(http.MethodPut)
}
