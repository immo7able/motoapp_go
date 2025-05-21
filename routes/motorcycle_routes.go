package routes

import (
	"github.com/gin-gonic/gin"
	"motorcycleApp/handler"
	"motorcycleApp/middleware"
)

func RegisterMotorcycleRoutes(r *gin.Engine, motorcycleHandler *handler.MotorcycleHandler, secretKey []byte) {
	motoGroup := r.Group("/motorcycles")
	r.Use(middleware.JWTAuthMiddleware(secretKey))
	{
		r.GET("/", motorcycleHandler.GetAllMotorcycles)
		r.GET("/motorcycles/:id", motorcycleHandler.ShowMotorcyclePage)
	}
	motoGroup.Use(middleware.JWTAuthSecuredMiddleware(secretKey))
	{
		motoGroup.GET("/add", motorcycleHandler.AddMotorcycle)
		motoGroup.POST("/add", motorcycleHandler.AddMotorcycle)
		motoGroup.GET("/my", motorcycleHandler.GetUserMotorcycles)
		motoGroup.POST("/delete/:id", motorcycleHandler.DeleteMotorcycle)
		motoGroup.GET("/edit/:id", motorcycleHandler.ShowEditPage)
		motoGroup.POST("/edit/:id", motorcycleHandler.EditMotorcycle)
	}
	commentGroup := r.Group("/comments")
	commentGroup.Use(middleware.JWTAuthSecuredMiddleware(secretKey))
	{
		commentGroup.POST("/:id", motorcycleHandler.AddComment)
		commentGroup.POST("/:id/delete", motorcycleHandler.DeleteComment)

	}
}
