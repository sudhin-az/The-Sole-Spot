package routes

import (
	"ecommerce_clean_architecture/pkg/api/handlers"
	"ecommerce_clean_architecture/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, cartHandler *handlers.CartHandler, orderHandler *handlers.OrderHandler, productHandler *handlers.ProductHandler) {
	router.POST("/usersignup", userHandler.UserSignup)
	router.POST("/verify-otp/:email", userHandler.VerifyOTP)
	router.POST("/resend-otp/:email", userHandler.ResendOTP)
	router.POST("/userlogin", userHandler.UserLogin)
	router.GET("/listproducts", userHandler.GetProducts)
	router.GET("/listcategory", userHandler.ListCategory)

	address := router.Group("/addresses")
	{
		address.Use(middleware.AuthMiddleware())
		address.POST("/addaddress", userHandler.AddAddress)
		address.PUT("/editaddress", userHandler.UpdateAddress)
		address.DELETE("/deleteaddress", userHandler.DeleteAddress)
		address.GET("/alladdress", userHandler.GetAllAddresses)
		address.Use(middleware.AuthMiddleware())
		address.GET("/userprofile", userHandler.UserProfile)
		address.PUT("/editprofile", userHandler.UpdateProfile)
		address.PUT("/forgotpassword", userHandler.ForgotPassword)
	}

	// product := router.Group("/products")
	// {
	// 	product.GET("/filtercategory", productHandler.FilterCategory)
	// 	product.GET("/searchproduct", productHandler.SearchProduct)
	// }

	cart := router.Group("/cart")
	{
		cart.Use(middleware.AuthMiddleware())
		cart.POST("/addtocart", cartHandler.AddToCart)
		cart.DELETE("/removefromcart", cartHandler.RemoveFromCart)
		cart.GET("/displaycart", cartHandler.DisplayCart)
	}

	// order := router.Group("/order")
	// {
	// 	order.POST("/place", orderHandler.PlaceOrder)
	// }
}
func AuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	router.GET("/google/login", authHandler.GoogleLogin)
	router.GET("/google/callback", authHandler.GoogleCallback)
}
