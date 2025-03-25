package routes

import (
	"ecommerce_clean_arch/pkg/api/handlers"
	"ecommerce_clean_arch/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.RouterGroup, adminHandler *handlers.AdminHandler, categoryHandler *handlers.CategoryHandler,
	productHandler *handlers.ProductHandler, couponHandler *handlers.CouponHandler) {
	router.POST("/adminsignup", adminHandler.SignUpHandler)
	router.POST("/adminlogin", adminHandler.LoginHandler)

	userdetails := router.Group("/users")
	{
		userdetails.Use(middleware.AdminMiddleware())
		userdetails.GET("/listofusers", adminHandler.GetUsers)
		userdetails.GET("/blockusers", adminHandler.BlockUser)
		userdetails.GET("/unblockusers", adminHandler.UnBlockUsers)
	}

	category := router.Group("/category")
	{
		category.Use(middleware.AdminMiddleware())
		category.POST("/addcategory", categoryHandler.AddCategory)
		category.PUT("/updatecategory", categoryHandler.UpdateCategory)
		category.DELETE("/deletecategory", categoryHandler.DeleteCategory)
	}

	product := router.Group("/product")
	{
		product.Use(middleware.AdminMiddleware())
		product.POST("/addproduct", productHandler.AddProduct)
		product.PUT("/updateproduct", productHandler.UpdateProduct)
		product.DELETE("/deleteproduct", productHandler.DeleteProduct)
	}

	orders := router.Group("/orders")
	{
		orders.Use(middleware.AdminMiddleware())
		orders.GET("/listorders", adminHandler.ListOrders)
		orders.PATCH("cancelorders", adminHandler.AdminCancelOrders)
		orders.PUT("/changeorderstatus", adminHandler.ChangeOrderStatus)
	}

	salesreportmanagement := router.Group("/salesreport")
	{
		salesreportmanagement.Use(middleware.AdminMiddleware())
		salesreportmanagement.GET("/get", adminHandler.SalesReport)
		salesreportmanagement.GET("/generateSalesreport", adminHandler.GenerateSalesReport)
		salesreportmanagement.GET("/BestSellingProduct", adminHandler.BestSellingProduct)
		salesreportmanagement.GET("/BestSellingCategory", adminHandler.BestSellingCategory)
	}

	coupon := router.Group("/coupons")
	{
		coupon.Use(middleware.AdminMiddleware())
		coupon.POST("/addcoupon", couponHandler.CreateNewCoupon)
		coupon.DELETE("/deletecoupon", couponHandler.MakeCouponInvalid)
		coupon.GET("/getallcoupon", couponHandler.GetAllCoupons)
	}
}
