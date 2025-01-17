package routes

import (
	"ecommerce_clean_architecture/pkg/api/handlers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.RouterGroup, adminHandler *handlers.AdminHandler, categoryHandler *handlers.CategoryHandler, productHandler *handlers.ProductHandler) {
	router.POST("/adminsignup", adminHandler.SignUpHandler)
	router.POST("/adminlogin", adminHandler.LoginHandler)

	userdetails := router.Group("/users")
	{
		userdetails.GET("/listofusers", adminHandler.GetUsers)
		userdetails.GET("/blockusers", adminHandler.BlockUser)
		userdetails.GET("/unblockusers", adminHandler.UnBlockUsers)
	}

	category := router.Group("/category")
	{
		category.POST("/addcategory", categoryHandler.AddCategory)
		category.PUT("/updatecategory", categoryHandler.UpdateCategory)
		category.DELETE("/deletecategory", categoryHandler.DeleteCategory)
	}

	product := router.Group("/product")
	{
		product.POST("/addproduct", productHandler.AddProduct)
		product.PUT("/updateproduct", productHandler.UpdateProduct)
		product.DELETE("/deleteproduct", productHandler.DeleteProduct)
	}
}
