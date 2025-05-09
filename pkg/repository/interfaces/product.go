package interfaces

import (
	"ecommerce_clean_arch/pkg/domain"
	"ecommerce_clean_arch/pkg/utils/models"
)

type ProductRepository interface {
	AddProduct(product models.AddProduct) (models.ProductResponse, error)
	UpdateProduct(products models.ProductResponse, productID int) (models.ProductResponse, error)
	DeleteProduct(productID int) error
	GetProductByID(productID int) (models.ProductResponse, error)
	GetProductsByCategory(categoryID string, sortBy string) ([]domain.Products, error)
}
