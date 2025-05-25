package interfaces

import (
	"ecommerce_clean_arch/pkg/domain"
	"ecommerce_clean_arch/pkg/utils/models"
)

type ProductUseCase interface {
	AddProduct(models.AddProduct) (models.ProductResponse, error)
	UpdateProduct(products models.ProductResponse, productID int) (models.ProductResponse, error)
	DeleteProduct(productID int) error
	SearchProduct(categoryID string, sortBy string) ([]domain.Products, error)
}
