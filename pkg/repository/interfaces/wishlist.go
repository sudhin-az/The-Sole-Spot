package interfaces

import "ecommerce_clean_arch/pkg/utils/models"

type WishlistRepository interface {
	AddToWishList(userID int, productID int) error
	GetWishList(userID int) ([]models.WishListResponse, error)
	RemoveFromWishList(userID int, productID int) error
	ProductExistInWishList(productID int, userID int) (bool, error)
	DoesProductExist(productID int) (bool, error)
}
