package main

import (
	"log"
	"os"

	"github.com/brkss/dextrace/internal/delivery"
	"github.com/brkss/dextrace/internal/domain"
	"github.com/brkss/dextrace/internal/infrastructure"
	"github.com/brkss/dextrace/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	user := domain.User{
		Email:    os.Getenv("USER_EMAIL"),
		Password: os.Getenv("USER_PASSWORD"),
	}
	userID := os.Getenv("USER_ID")

	// Initialize repositories
	sibionicRepo := infrastructure.NewSibionicRepository(os.Getenv("API_URL"))

	// Initialize use cases
	sibionicUseCase := usecase.NewSibionicUseCase(sibionicRepo, sibionicRepo)

	// Initialize handlers
	handler := delivery.NewGlucoseHandler(sibionicUseCase, userID, user)

	// Setup router
	r := gin.Default()
	r.GET("/data", handler.GetGlucoseData)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}