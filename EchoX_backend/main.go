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
	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.Record{},
		&models.Tag{},
		&models.MediaFile{},
		&models.PasswordResetToken{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize controllers
	userController := controllers.NewUserController(config.DB)
	recordController := controllers.NewRecordController(config.DB)
	tagController := controllers.NewTagController(config.DB)
	uploadController := controllers.NewUploadController(config.DB)
	passwordController := controllers.NewPasswordController(config.DB)
	discoveryController := controllers.NewDiscoveryController(config.DB)

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

	// 静态文件服务：提供上传的图片/视频访问（无需鉴权）
	r.Static("/uploads", "./uploads")

	// Setup all routes using the routes package
	routes.SetupRoutes(r, userController, recordController, tagController, uploadController, passwordController, discoveryController)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
