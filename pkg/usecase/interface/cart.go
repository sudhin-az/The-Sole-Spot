package interfaces

import "ecommerce_clean_architecture/pkg/utils/models"

type CartUseCase interface {
	AddToCart(userID int, productID int, quantity int) (models.CartResponse, error)
	DisplayCart(userID int) ([]models.Cart, error)
	RemoveProductFromCart(userID int, productID int) (models.CartResponse, error)
}
