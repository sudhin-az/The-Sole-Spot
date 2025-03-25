package interfaces

import "ecommerce_clean_arch/pkg/utils/models"

type WishlistUseCase interface {
	AddToWishList(productID int, userID int) error
	GetWishList(userID int) ([]models.WishListResponse, error)
	RemoveFromWishList(productID int, userID int) error
}
