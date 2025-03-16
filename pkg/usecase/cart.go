package usecase

import (
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/repository/interfaces"
	"ecommerce_clean_architecture/pkg/utils"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
)

type CartUseCase struct {
	cartRepository     repository.CartRepository
	productRepository  repository.ProductRepository
	categoryRepository repository.CategoryRepository
}

func NewCartUseCase(cartRepository repository.CartRepository, productrepository repository.ProductRepository, categoryRepository repository.CategoryRepository) *CartUseCase {
	return &CartUseCase{
		cartRepository:     cartRepository,
		productRepository:  productrepository,
		categoryRepository: categoryRepository,
	}
}

func (uc *CartUseCase) ValidateAddToCart(userID int, productID int, requestQty int) error {
	product, err := uc.productRepository.GetProductByID(productID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.Quantity < 0 {
		return errors.New("invalid quantity")
	}

	if requestQty > product.Stock {
		return errors.New("requested quantity exceeds available stock")
	}

	cartItem, err := uc.cartRepository.GetCartItem(userID, productID)
	if err != nil && err.Error() != "record not found" {
		return errors.New("error fetching cart details")
	}
	maxQtyPerPerson := 10
	totalQty := requestQty
	if cartItem != nil {
		totalQty += cartItem.Quantity
	}
	if totalQty > maxQtyPerPerson {
		return errors.New("exceeded max quantity allowed per person")
	}
	return nil
}

func (uc *CartUseCase) GetFilterProducts(showOutOfStock bool) ([]models.ProductResponse, error) {
	return uc.productRepository.GetAllProducts(showOutOfStock)
}
func (cu *CartUseCase) DisplayCart(userID int) ([]models.Cart, error) {
	cart, err := cu.cartRepository.DisplayCart(userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		return []models.Cart{}, nil
	}
	return cart, nil
}
func (cu *CartUseCase) AddToCart(userID int, productID int, quantity int) (models.CartResponse, error) {

	product, err := cu.productRepository.GetProductByID(productID)
	if err != nil {
		return models.CartResponse{}, errors.New("product not found")
	}

	if product.Quantity < 0 {
		return models.CartResponse{}, errors.New("invalid quantity")
	}

	if quantity > product.Quantity {
		return models.CartResponse{}, errors.New("insufficient quantity available")
	}

	if quantity > interfaces.MaxQuantity {
		return models.CartResponse{}, errors.New("quantity limit exceeded")
	}

	category, err := cu.categoryRepository.GetCategoryByID(product.Category_Id)
	if err != nil {
		return models.CartResponse{}, errors.New("category not found")
	}

	discount := category.CategoryDiscount
	discountedPrice := float64(product.OfferPrice) * (1 - (float64(discount) / 100))
	totalPrice := float64(quantity) * discountedPrice

	existingCartItem, _ := cu.cartRepository.GetCartItem(userID, productID)
	if existingCartItem != nil {
		existingCartItem.Quantity += quantity
		existingCartItem.TotalPrice = float64(existingCartItem.Quantity) * discountedPrice
		updatedCart, err := cu.cartRepository.UpdateCart(*existingCartItem)
		if err != nil {
			return models.CartResponse{}, err
		}
		return models.CartResponse{
			TotalPrice: updatedCart.TotalPrice,
			Cart:       []models.Cart{updatedCart},
		}, nil
	}

	newCartItem := models.Cart{
		UserID:     userID,
		ProductID:  productID,
		Quantity:   quantity,
		Price:      int(product.Price),
		OfferPrice: int(product.OfferPrice),
		CategoryDiscount: utils.RoundToTwoDecimalPlaces(product.OfferPrice - discountedPrice)*float64(quantity),
		TotalPrice: totalPrice,
	}

	addedCart, err := cu.cartRepository.AddToCart(newCartItem)
	if err != nil {
		return models.CartResponse{}, err
	}

	return models.CartResponse{
		TotalPrice: addedCart.TotalPrice,
		Cart:       []models.Cart{addedCart},
	}, nil
}

func (cu *CartUseCase) RemoveProductFromCart(userID int, productID int) (models.CartResponse, error) {
	product, err := cu.productRepository.GetProductByID(productID)
	if err != nil {
		return models.CartResponse{}, errors.New("product not found")
	}

	exists, err := cu.cartRepository.CheckProductInCart(userID, productID)
	if err != nil {
		return models.CartResponse{}, err
	}
	if !exists {
		return models.CartResponse{}, errors.New("product not found in the cart")
	}

	err = cu.cartRepository.RemoveProductFromCart(userID, productID, product.Price)
	if err != nil {
		return models.CartResponse{}, err
	}

	updatedCart, err := cu.cartRepository.DisplayCart(userID)
	if err != nil {
		return models.CartResponse{}, err
	}

	var totalPrice float64
	for _, item := range updatedCart {
		totalPrice += item.TotalPrice
	}

	return models.CartResponse{
		Cart:       updatedCart,
		TotalPrice: totalPrice,
	}, nil
}
