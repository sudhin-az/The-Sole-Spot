package interfaces

import "ecommerce_clean_architecture/pkg/utils/models"

type ProductUseCase interface {
	AddProduct(product models.AddProduct) (models.ProductResponse, error)
	UpdateProduct(products models.ProductResponse, productID int) (models.ProductResponse, error)
	DeleteProduct(productID int) error
}
