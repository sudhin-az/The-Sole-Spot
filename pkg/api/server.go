package api

import (
	"ecommerce_clean_architecture/pkg/api/handlers"
	"ecommerce_clean_architecture/pkg/api/routes"
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler,
	adminHandler *handlers.AdminHandler, categoryHandler *handlers.CategoryHandler, productHandler *handlers.ProductHandler,
	reviewHandler *handlers.ReviewHandler, cartHandler *handlers.CartHandler, orderHandler *handlers.OrderHandler,
	paymentHandler *handlers.PaymentHandler, walletHandler *handlers.WalletHandler, wishlistHandler *handlers.WishlistHandler, couponHandler *handlers.CouponHandler) *ServerHTTP {

	router := gin.New()

	//add swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	templatePath, _ := filepath.Abs("../templates/*")
	router.LoadHTMLGlob(templatePath)

	// Set up user routes
	userGroup := router.Group("/user")
	routes.UserRoutes(userGroup, userHandler, cartHandler, orderHandler, productHandler, reviewHandler,
		paymentHandler, walletHandler, wishlistHandler)

	authGroup := router.Group("/auth")
	routes.AuthRoutes(authGroup, authHandler)

	//Set up admin routes
	adminGroup := router.Group("/admin")
	routes.AdminRoutes(adminGroup, adminHandler, categoryHandler, productHandler, couponHandler)

	router.GET("/payment", paymentHandler.CreatePayment)
	router.POST("/payment/verify", paymentHandler.OnlinePaymentVerification)
	router.GET("/payment/success", paymentHandler.PaymentSuccess)

	return &ServerHTTP{engine: router}
}

// Start runs the HTTP server on a specified port
func (sh *ServerHTTP) Start(port string) {
	if err := sh.engine.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
