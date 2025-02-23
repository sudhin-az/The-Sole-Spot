package repository

import (
	"ecommerce_clean_architecture/pkg/domain"
	"errors"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(DB *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		DB: DB,
	}
}

func (cat *CategoryRepository) AddCategory(category domain.Category) (domain.Category, error) {
	var categoryResponse domain.Category

	err := cat.DB.Raw("INSERT INTO categories (category, description, category_discount) VALUES (?, ?, ?) RETURNING id, category, description, category_discount", category.Category, category.Description, category.CategoryDiscount).Scan(&categoryResponse).Error
	if err != nil {
		return domain.Category{}, err
	}
	return categoryResponse, nil
}

func (cat *CategoryRepository) UpdateCategory(category domain.Category, categoryID int) (domain.Category, error) {
	var updatedCategory domain.Category

	if err := cat.DB.Raw("SELECT id FROM categories WHERE id = ?", categoryID).Scan(&updatedCategory).Error; err != nil {
		return domain.Category{}, errors.New("category ID not found")
	}

	err := cat.DB.Raw(`
    UPDATE categories 
    SET category = ?, description = ?, category_discount = ?
    WHERE id = ? 
    RETURNING id, category, description, category_discount`,
		category.Category, category.Description, category.CategoryDiscount, categoryID).Scan(&updatedCategory).Error

	if err != nil {
		return domain.Category{}, err
	}

	if updatedCategory.ID == 0 {
		return domain.Category{}, errors.New("update failed; no rows affected")
	}

	return updatedCategory, nil
}

func (cat *CategoryRepository) DeleteCategory(categoryID int) error {
	var categories domain.Category

	result := cat.DB.Where("id = ?", categoryID).Delete(&categories)
	if result.RowsAffected < 1 {
		return errors.New("the ID does not exist")
	}
	return nil
}

func (p *CategoryRepository) GetCategoryByID(categoryID int) (domain.Category, error) {

	var categoryResponse domain.Category
	err := p.DB.Raw("SELECT * FROM categories WHERE id = ?", categoryID).Scan(&categoryResponse).Error
	if err != nil {
		return domain.Category{}, err
	}
	return categoryResponse, nil
}
