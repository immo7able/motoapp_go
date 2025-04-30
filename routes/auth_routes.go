package routes

import (
	"github.com/gin-gonic/gin"
	"motorcycleApp/handler"
)

func RegisterAuthRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
	auth := r.Group("/auth")
	{
		auth.GET("/register", authHandler.ShowRegisterPage)
		auth.POST("/register-form", authHandler.RegisterForm)

		auth.GET("/login", authHandler.ShowLoginPage)
		auth.POST("/login-form", authHandler.LoginForm)

		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}
