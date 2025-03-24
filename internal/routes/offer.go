package routes

import (
	"chase_deal/internal/handlers"
	"chase_deal/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterOfferRoutes(router *mux.Router, offerHandler *handlers.OfferHandler, shopOwnerMiddleware *middleware.ShopOwnerMiddleware) {
	// Public routes
	publicRouter := router.PathPrefix("/offers").Subrouter()
	publicRouter.HandleFunc("", offerHandler.ListOffers).Methods("GET")
	publicRouter.HandleFunc("/{id}", offerHandler.GetOffer).Methods("GET")

	// Shop owner routes
	ownerRouter := router.PathPrefix("/offers").Subrouter()
	ownerRouter.Use(middleware.RequireRole("vendor"))
	ownerRouter.HandleFunc("", offerHandler.CreateOffer).Methods("POST")
	ownerRouter.HandleFunc("/{id}", offerHandler.UpdateOffer).Methods("PUT")
	ownerRouter.HandleFunc("/{id}/feature", offerHandler.ToggleFeatured).Methods("PATCH")
}
