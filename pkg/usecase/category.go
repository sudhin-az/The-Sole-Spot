package usecase

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/repository"
)

type CategoryUseCase struct {
	CategoryRepository repository.CategoryRepository
}

func NewCategoryUseCase(usecase repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		CategoryRepository: usecase,
	}
}

func (cat *CategoryUseCase) AddCategory(category domain.Category) (domain.Category, error) {
	categoryResponse, err := cat.CategoryRepository.AddCategory(category)
	if err != nil {
		return domain.Category{}, err
	}
	return categoryResponse, nil
}

func (cat *CategoryUseCase) UpdateCategory(category domain.Category, categoryID int) (domain.Category, error) {
	updateCategory, err := cat.CategoryRepository.UpdateCategory(category, categoryID)
	if err != nil {
		return domain.Category{}, err
	}
	return updateCategory, nil
}

func (cat *CategoryUseCase) DeleteCategory(categoryID int) error {

	err := cat.CategoryRepository.DeleteCategory(categoryID)
	if err != nil {
		return err
	}
	return nil
}
