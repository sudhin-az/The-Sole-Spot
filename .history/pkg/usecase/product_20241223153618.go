package usecase

import (
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
)

type ProductUseCase struct {
	ProductRepository repository.ProductRepository
}

func NewProductUseCase(usecase repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		ProductRepository: usecase,
	}
}

func (p *ProductUseCase) AddProduct(product models.AddProduct) (models.ProductResponse, error) {
	if product.Price < 0 || product.Quantity < 0 {
		return models.ProductResponse{}, errors.New("invalid quantity or price")
	}
	products, err := p.ProductRepository.AddProduct(product)
	if err != nil {
		return models.ProductResponse{}, err
	}
	productResponse := models.ProductResponse{
		ID:          products.ID,
		Category_Id: products.CategoryID,
		Name:        products.Name,
		Price:       products.Price,
		Quantity:    products.Quantity,
	}

	return productResponse, nil
}

func (p *ProductUseCase) UpdateProduct(products models.ProductResponse, productID int) (models.ProductResponse, error) {

	if products.Price < 0 || products.Quantity < 0 {
		return models.ProductResponse{}, errors.New("invalid quantity or price")
	}

	updateProduct, err := p.ProductRepository.UpdateProduct(products, productID)
	if err != nil {
		return models.ProductResponse{}, err
	}
	return updateProduct, nil
}

func (p *ProductUseCase) DeleteProduct(productID int) error {
	err := p.ProductRepository.DeleteProduct(productID)
	if err != nil {
		return err
	}
	return nil
}
