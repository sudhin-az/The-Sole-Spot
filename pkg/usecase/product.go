package usecase

import (
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils/models"
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

	products, err := p.ProductRepository.AddProduct(product)
	if err != nil {
		return models.ProductResponse{}, err
	}
	return products, nil
}

func (p *ProductUseCase) UpdateProduct(products models.ProductResponse, productID int) (models.ProductResponse, error) {

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
