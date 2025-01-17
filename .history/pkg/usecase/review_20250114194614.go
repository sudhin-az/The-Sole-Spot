package usecase

import (
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils/models"
)

type ReviewUseCase struct {
	repo repository.ReviewRepository
}

func NewReviewUseCase(repository repository.ReviewRepository) *ReviewUseCase {
	return &ReviewUseCase{repo: repository}
}

func (r *ReviewUseCase) AddReview(userID string, productID string, Rating int, Comment string) error {
	err := r.repo.AddReview(userID, productID, Rating, Comment)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReviewUseCase) GetReviewsByProductID(productID string) ([]models.ReviewResponse, error) {

	reviews, err := r.repo.GetReviewsByProductID(productID)
	if err != nil {
		return []models.ReviewResponse{}, err
	}
	return reviews, nil
}

func (r *ReviewUseCase) DeleteReview(reviewID string) error {
	err := r.repo.DeleteReview(reviewID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReviewUseCase) GetAverageRating(productID string) (float64, error) {
	return r.repo.GetAverageRating(productID)
}
