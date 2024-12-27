package main

import (
	"go-ecommerce-api/database"
	"go-ecommerce-api/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-ecommerce-api/docs"
)

// @title E-Commerce API
// @version 1.0
// @description This is the API documentation for the E-Commerce application.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://example.com/contact
// @contact.email support@example.com

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found. Using system environment variables.")
	}

	// Connect to the database
	database.ConnectToDatabase()

	// Set up Gin router and routes
	gin.SetMode(gin.DebugMode)
	router := routes.SetupRoutes()
	router.SetTrustedProxies(nil)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Starting server on port 8080...")
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
