package routes

import (
	"github.com/gin-gonic/gin"
	"motorcycleApp/handler"
	"motorcycleApp/middleware"
)

func RegisterAdminRoutes(r *gin.Engine, adminHandler *handler.AdminHandler, secretKey []byte) {
	admin := r.Group("/admin", middleware.AdminOnly(secretKey))
	{
		admin.GET("/brands", adminHandler.CreateBrand)
		admin.POST("/brands", adminHandler.CreateBrand)

		admin.GET("/models", adminHandler.CreateModel)
		admin.POST("/models", adminHandler.CreateModel)
	}
}
