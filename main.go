package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"motorcycleApp/config"
	"motorcycleApp/handler"
	"motorcycleApp/routes"
	"motorcycleApp/service"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	cfg := config.NewConfig()
	config.ConnectDatabase(cfg.DatabaseUrl)

	db := config.DB

	authService := &service.AuthService{
		DB:     db,
		JWTKey: []byte(cfg.JWT.Secret),
	}

	newValidator := validator.New()
	authHandler := &handler.AuthHandler{
		AuthService: authService,
		Validator:   newValidator,
	}

	routes.RegisterAuthRoutes(router, authHandler)

	err := router.Run(cfg.ServerAddress)
	if err != nil {
		return
	}
}
