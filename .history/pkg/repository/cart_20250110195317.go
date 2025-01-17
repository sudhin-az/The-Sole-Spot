package repository

import (
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type CartRepository struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{
		DB: db,
	}
}

func (car *CartRepository) DisplayCart(userID int) ([]models.Cart, error) {
	var cartResponse []models.Cart
	err := car.DB.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&cartResponse).Error
	if err != nil {
		return nil, err
	}
	return cartResponse, nil
}

func (r *CartRepository) GetCartItem(userID int, productID int) (*models.Cart, error) {
	var cartItem models.Cart
	fmt.Println("Debug: Running query with userID=", userID, "productID=", productID)
	if err := r.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&cartItem).Error; err != nil {
		fmt.Println("Debug: Query error -", err)
		return nil, err
	}
	fmt.Println("Debug: Query successful, cartItem=", cartItem)
	return &cartItem, nil
}
func (car *CartRepository) AddToCart(cart models.Cart) (models.Cart, error) {
	err := car.DB.Create(&cart).Error
	if err != nil {
		return models.Cart{}, err
	}
	return cart, nil
}

func (car *CartRepository) UpdateCart(cart models.Cart) (models.Cart, error) {
	err := car.DB.Save(&cart).Error
	if err != nil {
		return models.Cart{}, err
	}
	return cart, nil
}
func (cr *CartRepository) CheckProductInCart(userID int, productID int) (bool, error) {
	var count int
	err := cr.DB.Raw(`SELECT COUNT(*) FROM carts WHERE user_id = $1 AND product_id = $2`, userID, productID).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (car *CartRepository) RemoveProductFromCart(userID int, productID int) error {
	// Find one cart entry and delete it
	var cartItem models.Cart
	err := car.DB.Where("user_id = ? AND product_id = ?", userID, productID).
		Order("id").
		First(&cartItem).Error
	if err != nil {
		return err
	}

	// Delete the found entry
	return car.DB.Delete(&cartItem).Error
}

func (car *CartRepository) GetAllItemsFromCart(userID int) ([]models.Cart, error) {
	var count int
	var cartResponse []models.Cart

	err := car.DB.Raw("select count(*) from carts where user_id = ? AND deleted_at IS NULL", userID).Scan(&count).Error
	if err != nil {
		return []models.Cart{}, err
	}

	if count == 0 {
		return []models.Cart{}, nil
	}
	err = car.DB.Raw("select carts.user_id,users.name as user_name,carts.product_id,products.name as name,carts.quantity,carts.total_price from carts inner join users on carts.user_id = users.id inner join products on carts.product_id = products.id where user_id = ? AND carts.deleted_at IS NULL", userID).First(&cartResponse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if len(cartResponse) == 0 {
				return []models.Cart{}, nil
			}
			return []models.Cart{}, err
		}
		return []models.Cart{}, err
	}

	return cartResponse, nil
}
