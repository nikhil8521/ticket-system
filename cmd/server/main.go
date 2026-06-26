package main

import (
	"log"
	"os"
	"ticket-system/controllers"
	"ticket-system/database"
	"ticket-system/repository"
	"ticket-system/routes"
	"ticket-system/services"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	database.InitDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.DB)
	ticketRepo := repository.NewTicketRepository(database.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	ticketService := services.NewTicketService(ticketRepo)

	// Initialize controllers
	authCtrl := controllers.NewAuthController(authService)
	ticketCtrl := controllers.NewTicketController(ticketService)

	// Setup router
	r := routes.SetupRouter(authCtrl, ticketCtrl)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
