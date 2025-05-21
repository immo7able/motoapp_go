package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"motorcycleApp/config"
	"motorcycleApp/handler"
	"motorcycleApp/routes"
	"motorcycleApp/service"
)

func SetupApp() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*/*.html")
	router.Static("/css", "./static/css")
	router.Static("/uploads", "./uploads")
	router.Static("/js", "./static/js")

	cfg := config.NewConfig()
	config.ConnectDatabase(cfg.DatabaseUrl)
	db := config.DB

	adminService := &service.AdminService{
		DB: db,
	}

	adminHandler := &handler.AdminHandler{
		Service:   adminService,
		Validator: validator.New(),
	}

	authService := &service.AuthService{
		DB:     db,
		JWTKey: []byte(cfg.JWT.Secret),
	}

	authHandler := &handler.AuthHandler{
		AuthService: authService,
		Validator:   validator.New(),
	}

	motorcycleService := &service.MotorcycleService{
		DB:     db,
		JWTKey: []byte(cfg.JWT.Secret),
	}

	motorcycleHandler := &handler.MotorcycleHandler{
		Service:      motorcycleService,
		Validator:    validator.New(),
		AdminService: adminService,
	}

	routes.RegisterAuthRoutes(router, authHandler, []byte(cfg.JWT.Secret))
	routes.RegisterMotorcycleRoutes(router, motorcycleHandler, []byte(cfg.JWT.Secret))
	routes.RegisterAdminRoutes(router, adminHandler, []byte(cfg.JWT.Secret))

	return router
}
