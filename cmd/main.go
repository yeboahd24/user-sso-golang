package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yeboahd24/user-sso/config"
	"github.com/yeboahd24/user-sso/database"
	"github.com/yeboahd24/user-sso/handler"
	"github.com/yeboahd24/user-sso/repository"
	"github.com/yeboahd24/user-sso/route"
	"github.com/yeboahd24/user-sso/service"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, cfg)

	tmplHandler, err := handler.NewTemplateHandler()
	if err != nil {
		log.Fatal("Failed to initialize template handler:", err)
	}

	// Setup router
	r := gin.Default()
	route.SetupRoutes(r, authHandler, tmplHandler)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
