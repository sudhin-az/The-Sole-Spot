package interfaces

import "ecommerce_clean_arch/pkg/domain"

type CategoryRepository interface {
	AddCategory(category domain.Category) (domain.Category, error)
	UpdateCategory(category domain.Category, categoryID int) (domain.Category, error)
	DeleteCategory(categoryID int) error
	GetCategoryByID(categoryID int) (domain.Category, error)
}
