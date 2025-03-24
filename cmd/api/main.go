package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"

	"chase_deal/internal/database"
	"chase_deal/internal/handlers"
	"chase_deal/internal/middleware"
	"chase_deal/internal/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {

	hash, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	fmt.Println(string(hash))

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found")
	}

	// Initialize database connection pool
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Database closure error: %v", err)
		}
	}()

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Database migration error: %v", err)
	}

	// Initialize handlers with database connection
	authHandler := &handlers.AuthHandler{DB: db.SQL}
	adminHandler := &handlers.AdminHandler{DB: db.SQL}
	adminShopHandler := &handlers.AdminShopHandler{DB: db.SQL}

	// Create router
	router := mux.NewRouter()

	// Global middleware
	router.Use(middleware.Logging)
	router.Use(middleware.JSONContentType)

	// Public routes
	routes.RegisterAuthRoutes(router, authHandler)
	//routes.RegisterAdminRoutes(router, adminHandler)
	routes.RegisterAdminRoutes(router, adminHandler, adminShopHandler)

	// Protected routes
	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(middleware.JWTAuth)
	//routes.RegisterAuthRoutes(protectedRouter, bookHandler)

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
