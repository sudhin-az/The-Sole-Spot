package interfaces

import "ecommerce_clean_architecture/pkg/utils/models"

const MaxQuantity = 5

type CartRepository interface {
	DisplayCart(userID int) ([]models.Cart, error)
	GetCartItem(userID int, productID int) (models.Cart, error)
	AddToCart(cart models.Cart) (models.Cart, error)
	UpdateCart(cart models.Cart) (models.Cart, error)
	CheckProductInCart(userID int, productID int) (bool, error)
	RemoveProductFromCart(userID int, productID int) error
}
