package routes

import (
	"echox/internal/controllers"
	"echox/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	userController *controllers.UserController,
	recordController *controllers.RecordController,
	tagController *controllers.TagController,
	uploadController *controllers.UploadController,
	passwordController *controllers.PasswordController,
	discoveryController *controllers.DiscoveryController,
) {
	api := r.Group("/api/v1")

	// 公开路由（无需 JWT）
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", userController.Register)
		authGroup.POST("/login", userController.Login)
		authGroup.POST("/forgot-password", passwordController.ForgotPassword)
		authGroup.POST("/reset-password", passwordController.ResetPassword)
	}

	// 需要认证的路由
	auth := api.Group("")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		// 用户相关
		auth.GET("/users/profile", userController.GetProfile)
		auth.PUT("/users/profile", userController.UpdateProfile)

		// 记录相关
		auth.POST("/records", recordController.CreateRecord)
		auth.GET("/records", recordController.GetRecords)
		auth.GET("/records/:id", recordController.GetRecord)
		auth.PUT("/records/:id", recordController.UpdateRecord)
		auth.DELETE("/records/:id", recordController.DeleteRecord)

		// 标签相关
		auth.POST("/tags", tagController.CreateTag)
		auth.GET("/tags", tagController.GetTags)
		auth.PUT("/tags/:id", tagController.UpdateTag)
		auth.DELETE("/tags/:id", tagController.DeleteTag)

		// 文件上传与媒体
		auth.POST("/upload", uploadController.Upload)
		auth.GET("/media/:id", uploadController.GetMedia)
		auth.DELETE("/media/:id", uploadController.DeleteMedia)

		// 搜索与时间轴
		auth.GET("/search", discoveryController.Search)
		auth.GET("/timeline", discoveryController.Timeline)
	}
}
