package repository

import (
	"ecommerce_clean_arch/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type WishlistRepository struct {
	DB *gorm.DB
}

func NewWishlistRepository(DB *gorm.DB) *WishlistRepository {
	return &WishlistRepository{DB}
}
func (w *WishlistRepository) AddToWishList(userID int, productID int) error {
	query := "INSERT INTO wishlists(user_id, product_id) VALUES(?, ?)"
	if err := w.DB.Exec(query, userID, productID).Error; err != nil {
		return errors.New("encountered an issue while inserting into wishlist")
	}
	return nil
}

func (w *WishlistRepository) GetWishList(userID int) ([]models.WishListResponse, error) {
	var wishlist []models.WishListResponse
	query := "SELECT products.id as product_id, products.name as product_name, products.price as product_price FROM products INNER JOIN wishlists ON products.id = wishlists.product_id WHERE wishlists.user_id = ?"
	err := w.DB.Raw(query, userID).Scan(&wishlist).Error
	if err != nil {
		return []models.WishListResponse{}, errors.New("encountered an issue while fetching products from wishlist")
	}
	return wishlist, nil
}

func (w *WishlistRepository) RemoveFromWishList(userID int, productID int) error {
	query := "DELETE FROM wishlists WHERE product_id = ? AND user_id = ?"
	result := w.DB.Exec(query, productID, userID)
	if result.Error != nil {
		return errors.New("encountered an issue while deleting from wishlist")
	}
	if result.RowsAffected == 0 {
		return errors.New("no product was deleted, maybe it didn't exist")
	}
	return nil
}
func (w *WishlistRepository) ProductExistInWishList(productID int, userID int) (bool, error) {
	var count int64

	query := "SELECT COUNT(*) FROM wishlists WHERE product_id = ? AND user_id = ?"
	if err := w.DB.Raw(query, productID, userID).Scan(&count).Error; err != nil {
		return false, errors.New("error while checking wishlist")
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}
func (w *WishlistRepository) DoesProductExist(productID int) (bool, error) {
	var count int
	err := w.DB.Raw("select count(*) from products where id = ?", productID).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
