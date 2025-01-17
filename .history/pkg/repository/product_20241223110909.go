package repository

import (
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) *ProductRepository {
	return &ProductRepository{
		DB: DB,
	}
}

func (p *ProductRepository) AddProduct(addProduct models.AddProduct) (models.ProductResponse, error) {
	var products models.ProductResponse

	// Check if the category exists
	var categoryExists bool
	err := p.DB.Raw("SELECT EXISTS (SELECT 1 FROM categories WHERE id = ?)", addProduct.CategoryID).Scan(&categoryExists).Error
	if err != nil {
		return models.ProductResponse{}, fmt.Errorf("error checking category existence: %w", err)
	}

	if !categoryExists {
		return models.ProductResponse{}, errors.New("category does not exist")
	}

	// Insert the product
	err = p.DB.Raw("INSERT INTO products (category_id, name, quantity, stock, price) VALUES (?, ?, ?, ?, ?) RETURNING category_id, name, quantity, stock, price",
		addProduct.CategoryID, addProduct.Name, addProduct.Quantity, addProduct.Stock, addProduct.Price).Scan(&products).Error
	if err != nil {
		return models.ProductResponse{}, fmt.Errorf("error adding product: %w", err)
	}

	return products, nil
}

func (p *ProductRepository) UpdateProduct(products models.ProductResponse, productID int) (models.ProductResponse, error) {
	var productResponse models.ProductResponse

	err := p.DB.Raw("UPDATE products SET category_id = ?, name = ?, quantity = ?, stock = ?, price = ? WHERE id = ? RETURNING id, category_id, name, quantity, stock, price",
		products.CategoryID, products.Name, products.Quantity, products.Stock, products.Price, productID).Scan(&productResponse).Error
	if err != nil {
		return models.ProductResponse{}, fmt.Errorf("error updating product: %w", err)
	}

	return productResponse, nil
}

func (p *ProductRepository) DeleteProduct(productID int) error {
	err := p.DB.Exec("DELETE FROM products WHERE id = ?", productID)
	if err.RowsAffected < 1 {
		return errors.New("the id is not existing")
	}
	return nil
}
