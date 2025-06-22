package api

import (
	"ecommerce_clean_arch/pkg/api/handlers"
	"ecommerce_clean_arch/pkg/api/routes"
	"log"

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

	if userHandler == nil || authHandler == nil || adminHandler == nil || categoryHandler == nil || productHandler == nil ||
		reviewHandler == nil || cartHandler == nil || orderHandler == nil || paymentHandler == nil || walletHandler == nil ||
		wishlistHandler == nil || couponHandler == nil {
		log.Fatal("One or more handlers are nil")
	}

	router := gin.New()
	if router == nil {
		log.Fatal("Failed to initialize gin.Engine")
	}

	//add swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.LoadHTMLGlob("templates/*")

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

	log.Println("ServerHTTP initialized successfully")
	return &ServerHTTP{engine: router}
}

// Start runs the HTTP server on a specified port
func (sh *ServerHTTP) Start(port string) {
	if sh.engine == nil {
		log.Fatal("Engine is nil. Cannot start")
	}
	if err := sh.engine.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
