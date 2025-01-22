package interfaces

import "ecommerce_clean_architecture/pkg/utils/models"

type CartUseCase interface {
	GetFilterProducts(showOutOfStock bool) ([]models.ProductResponse, error)
	DisplayCart(userID int) ([]models.Cart, error)
	AddToCart(userID int, productID int, quantity int) (models.CartResponse, error)
	RemoveProductFromCart(userID int, productID int, Quantity int) (models.CartResponse, error)
}
