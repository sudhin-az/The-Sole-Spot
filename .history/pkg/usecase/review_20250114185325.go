package usecase

import "ecommerce_clean_architecture/pkg/repository"

type ReviewUseCase struct {
	repo repository.ReviewRepository
}
