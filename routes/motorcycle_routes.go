package routes

import (
	"github.com/gin-gonic/gin"
	"motorcycleApp/handler"
	"motorcycleApp/middleware"
)

func RegisterMotorcycleRoutes(r *gin.Engine, motorcycleHandler *handler.MotorcycleHandler, secretKey []byte) {
	motoGroup := r.Group("/motorcycles")
	motoGroup.Use(middleware.JWTAuthMiddleware(secretKey)).GET("/", motorcycleHandler.GetAllMotorcycles)
	motoGroup.Use(middleware.JWTAuthSecuredMiddleware(secretKey))
	{
		motoGroup.GET("/add", motorcycleHandler.AddMotorcycle)
		motoGroup.POST("/add", motorcycleHandler.AddMotorcycle)
		motoGroup.GET("/my", motorcycleHandler.GetUserMotorcycles)
		motoGroup.POST("/delete/:id", motorcycleHandler.DeleteMotorcycle)
	}
}
