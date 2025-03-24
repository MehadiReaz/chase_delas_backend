package routes

import (
	"net/http"

	"chase_deal/internal/handlers"
	"chase_deal/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterShopRoutes(
	router *mux.Router,
	adminShopHandler *handlers.AdminShopHandler,
	shopHandler *handlers.ShopHandler,
) {
	// Admin routes
	adminRouter := router.PathPrefix("/admin/shops").Subrouter()
	adminRouter.Use(middleware.JWTAuth)
	adminRouter.Use(middleware.RequireRole("superadmin"))

	adminRouter.HandleFunc("", adminShopHandler.ListShops).Methods(http.MethodGet)
	adminRouter.HandleFunc("", adminShopHandler.CreateShop).Methods(http.MethodPost)
	adminRouter.HandleFunc("/{id}/priority", adminShopHandler.UpdateShopPriority).Methods(http.MethodPut)

	// Shop owner routes
	shopRouter := router.PathPrefix("/shops").Subrouter()
	shopRouter.Use(middleware.JWTAuth)
	shopRouter.Use(middleware.RequireRole("vendor"))

	shopRouter.HandleFunc("/{id}", shopHandler.GetShopDetails).Methods(http.MethodGet)
	shopRouter.HandleFunc("/{id}", shopHandler.UpdateShop).Methods(http.MethodPut)
}
