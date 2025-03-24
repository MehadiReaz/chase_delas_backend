package routes

import (
	"net/http"

	"chase_deal/internal/handlers"
	"chase_deal/internal/middleware"
	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, userHandler *handlers.UserHandler) {
	userRouter := router.PathPrefix("/user").Subrouter()

	userRouter.Use(middleware.JWTAuth)
	userRouter.Use(middleware.RoleMiddleware("user", "vendor", "admin"))

	userRouter.HandleFunc("/profile", userHandler.GetProfile).Methods(http.MethodGet)
	userRouter.HandleFunc("/profile", userHandler.UpdateProfile).Methods(http.MethodPut)
}
