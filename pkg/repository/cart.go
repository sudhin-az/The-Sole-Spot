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
	var cartItems []models.Cart
	err := car.DB.Where("user_id = ?", userID).Find(&cartItems).Error
	if err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (r *CartRepository) GetCartItem(userID int, productID int) (*models.Cart, error) {
	var cartItem models.Cart
	fmt.Println("Running query with userID=", userID, "productID=", productID)
	if err := r.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&cartItem).Error; err != nil {
		fmt.Println("Query error -", err)
		return nil, err
	}
	fmt.Println("Query successful, cartItem=", cartItem)
	return &cartItem, nil
}
func (car *CartRepository) AddToCart(cartItem models.Cart) (models.Cart, error) {
	var existingCartItem models.Cart

	err := car.DB.Where("user_id = ? AND product_id = ?", cartItem.UserID, cartItem.ProductID).
		First(&existingCartItem).Error

	if err == nil {

		existingCartItem.Quantity += cartItem.Quantity
		existingCartItem.TotalPrice += cartItem.TotalPrice
		err = car.DB.Save(&existingCartItem).Error
		return existingCartItem, err
	}

	err = car.DB.Create(&cartItem).Error
	return cartItem, err
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

func (car *CartRepository) RemoveProductFromCart(userID int, productID int, price float64) error {
	var cartItem models.Cart

	err := car.DB.Where("user_id = ? AND product_id = ?", userID, productID).
		First(&cartItem).Error
	if err != nil {
		return errors.New("product not found in the cart")
	}

	if cartItem.Quantity > 1 {
		cartItem.Quantity--
		cartItem.TotalPrice -= float64(cartItem.OfferPrice)
		return car.DB.Save(&cartItem).Error
	}

	return car.DB.Delete(&cartItem).Error
}
func (car *CartRepository) RemoveFromCart(userID int, productID int) error {
	var cart models.Cart
	err := car.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error
	if err != nil {
		return errors.New("product not found in the cart")
	}
	return car.DB.Delete(&cart).Error
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
