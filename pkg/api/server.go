package api

import (
	"ecommerce_clean_architecture/pkg/api/handlers"
	"ecommerce_clean_architecture/pkg/api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

// NewServerHTTP initializes a new HTTP server with the necessary routes and handlers.
func NewServerHTTP(userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler,
	adminHandler *handlers.AdminHandler, categoryHandler *handlers.CategoryHandler, productHandler *handlers.ProductHandler) *ServerHTTP {
	router := gin.New()

	// Set up user routes
	userGroup := router.Group("/user")
	routes.UserRoutes(userGroup, userHandler)

	authGroup := router.Group("/auth")
	routes.AuthRoutes(authGroup, authHandler)

	//Set up admin routes
	adminGroup := router.Group("/admin")
	routes.AdminRoutes(adminGroup, adminHandler, categoryHandler, productHandler)

	return &ServerHTTP{engine: router}
}

// Start runs the HTTP server on a specified port
func (sh *ServerHTTP) Start(port string) {
	if err := sh.engine.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
