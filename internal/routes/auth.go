package routes

import (
	"net/http"

	"chase_deal/internal/handlers"
	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router, authHandler *handlers.AuthHandler) {
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/update-credentials", authHandler.UpdateCredentials).Methods("PUT")
}
