package routes

import (
	"github.com/gin-gonic/gin"
	"motorcycleApp/handler"
	"motorcycleApp/middleware"
)

func RegisterAdminRoutes(r *gin.Engine, adminHandler *handler.AdminHandler) {
	admin := r.Group("/admin", middleware.AdminOnly())
	{
		admin.GET("/brands", adminHandler.BrandsPage)
		admin.POST("/brands", adminHandler.BrandsPage)

		admin.GET("/models", adminHandler.ModelsPage)
		admin.POST("/models", adminHandler.ModelsPage)
	}
}
