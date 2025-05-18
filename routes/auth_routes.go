package routes

import (
	"github.com/gin-gonic/gin"
	"motorcycleApp/handler"
	"motorcycleApp/middleware"
)

func RegisterAuthRoutes(r *gin.Engine, authHandler *handler.AuthHandler, secretKey []byte) {
	auth := r.Group("/auth")
	auth.Use(middleware.JWTAuthRedirectMiddleware(secretKey))
	{
		auth.GET("/register", authHandler.Register)
		auth.POST("/register", authHandler.Register)

		auth.GET("/login", authHandler.Login)
		auth.POST("/login", authHandler.Login)

		auth.GET("/logout", authHandler.Logout)
	}
}
