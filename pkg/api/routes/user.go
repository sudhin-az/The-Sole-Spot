package routes

import (
	"ecommerce_clean_architecture/pkg/api/handlers"
	"ecommerce_clean_architecture/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, cartHandler *handlers.CartHandler,
	orderHandler *handlers.OrderHandler, productHandler *handlers.ProductHandler, reviewHandler *handlers.ReviewHandler,
	paymentHandler *handlers.PaymentHandler, walletHandler *handlers.WalletHandler, wishlistHandler *handlers.WishlistHandler) {
	router.Use(gin.Logger(), gin.Recovery())
	router.POST("/usersignup", userHandler.UserSignup)
	router.POST("/verify-otp/:email", userHandler.VerifyOTP)
	router.POST("/resend-otp/:email", userHandler.ResendOTP)
	router.POST("/userlogin", userHandler.UserLogin)
	router.GET("/listproducts", userHandler.GetProducts)
	router.GET("/listcategory", userHandler.ListCategory)

	//addresses
	address := router.Group("/addresses")
	{
		address.POST("/forgotpassword", userHandler.ForgotPassword)
		address.Use(middleware.AuthMiddleware())
		address.POST("/addaddress", userHandler.AddAddress)
		address.PUT("/editaddress", userHandler.UpdateAddress)
		address.DELETE("/deleteaddress", userHandler.DeleteAddress)
		address.GET("/alladdress", userHandler.GetAllAddresses)
		//userProfile
		address.Use(middleware.AuthMiddleware())
		address.GET("/userprofile", userHandler.UserProfile)
		address.PUT("/editprofile", userHandler.UpdateProfile)
		address.PUT("/changepassword", userHandler.ChangePassword)
	}

	//products
	product := router.Group("/products")
	// {
	// 	product.GET("/filtercategory", productHandler.FilterCategory)
	product.GET("/searchproduct", productHandler.SearchProduct)
	// }

	//cart
	cart := router.Group("/cart")
	{
		cart.Use(middleware.AuthMiddleware())
		cart.POST("/addtocart", cartHandler.AddToCart)
		cart.DELETE("/removefromcart", cartHandler.RemoveFromCart)
		cart.GET("/displaycart", cartHandler.DisplayCart)
	}

	//order
	order := router.Group("/order")
	{
		order.Use(middleware.AuthMiddleware())
		order.POST("/placeorder", orderHandler.OrderItemsFromCart)
		order.GET("/vieworders", orderHandler.ViewOrders)
		order.PUT("/cancelorders", orderHandler.CancelOrders)
		order.PUT("/cancelOrderItem", orderHandler.CancelOrderItem)
		order.PUT("/returnorder", orderHandler.ReturnUserOrder)
	}

	router.GET("/payment", paymentHandler.CreatePayment)
	router.POST("/payment/verify", paymentHandler.OnlinePaymentVerification)
	router.GET("/payment/success", paymentHandler.PaymentSuccess)

	//Wallet
	wallet := router.Group("/wallet")
	{
		wallet.Use(middleware.AuthMiddleware())
		wallet.GET("/getwallet", walletHandler.ViewWallet)
		wallet.GET("/wallethistory", walletHandler.GetWalletTransaction)
	}

	//Wishlist
	wishlist := router.Group("wishlist")
	{
		wishlist.Use(middleware.AuthMiddleware())
		wishlist.POST("/addtowishlist", wishlistHandler.AddToWishList)
		wishlist.GET("/getwishlist", wishlistHandler.GetWishList)
		wishlist.DELETE("/removewishlist", wishlistHandler.RemoveFromWishList)
	}

	//review
	Review := router.Group("/review")
	Review.Use(middleware.AuthMiddleware())
	Review.POST("/addreview", reviewHandler.AddReview)
	Review.GET("/getreviewsbyproduct", reviewHandler.GetReviewsByProductID)
	Review.DELETE("/removereview", reviewHandler.DeleteReview)
	Review.GET("/getaveragerating", reviewHandler.GetAverageRating)
}

func AuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	router.GET("/google/login", authHandler.GoogleLogin)
	router.GET("/google/callback", authHandler.GoogleCallback)
}
