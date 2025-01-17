package usecase

import "ecommerce_clean_architecture/pkg/repository"

type ReviewUseCase struct {
	repo repository.ReviewRepository
}

func NewReviewUseCase(repository repository.ReviewRepository) *ReviewUseCase {
	return &ReviewUseCase{repo: repository}
}

func (r *ReviewUseCase) AddReview(userID string, productID string, Rating int, Comment string) error {

}
