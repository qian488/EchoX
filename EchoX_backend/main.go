package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"echox/internal/config"
	"echox/internal/controllers"
	"echox/internal/models"
	"echox/internal/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database
	config.InitDB()

	// Auto migrate database schema
	if err := config.DB.AutoMigrate(&models.User{}, &models.Record{}, &models.Tag{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize controllers
	userController := controllers.NewUserController(config.DB)
	recordController := controllers.NewRecordController(config.DB)
	tagController := controllers.NewTagController(config.DB)

	// Initialize router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup all routes using the routes package
	routes.SetupRoutes(r, userController, recordController, tagController)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
