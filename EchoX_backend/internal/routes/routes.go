package routes

import (
	"echox/internal/controllers"
	"echox/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userController *controllers.UserController, recordController *controllers.RecordController, tagController *controllers.TagController) {
	api := r.Group("/api/v1")

	// 公开路由
	authGroup := api.Group("/auth") // 创建 /api/v1/auth 路由组
	{
		authGroup.POST("/register", userController.Register) // 路径变为 /api/v1/auth/register
		authGroup.POST("/login", userController.Login)       // 路径变为 /api/v1/auth/login
	}

	// 需要认证的路由
	auth := api.Use(middleware.JWTAuthMiddleware()) // 认证路由组仍然在 /api/v1 下，但现在可以考虑是否也放到 /auth 组内
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
	}
}
