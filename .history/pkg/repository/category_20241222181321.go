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

	err := cat.DB.Raw("INSERT INTO categories (category, description) VALUES (?, ?) RETURNING category, description", category.Category, category.Description).Scan(&categoryResponse).Error
	if err != nil {
		return domain.Category{}, err
	}
	return categoryResponse, nil
}

func (cat *CategoryRepository) UpdateCategory(category domain.Category, categoryID int) (domain.Category, error) {
	var updatedCategory domain.Category
	err := cat.DB.Raw("UPDATE categories SET category = ?, description = ? WHERE id = ? RETURNING id, category, description", category.Category, category.Description, categoryID).Scan(&updatedCategory).Error
	if err != nil {
		return domain.Category{}, err
	}
	return updatedCategory, nil
}

func (cat *CategoryRepository) DeleteCategory(categoryID int) error {

	err := cat.DB.Exec("DELETE FROM categories WHERE id = ?", categoryID)
	if err.RowsAffected < 1 {
		return errors.New("the id is not existing")
	}
	return nil
}
