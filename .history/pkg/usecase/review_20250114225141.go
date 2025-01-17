package usecase

import (
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
)

type ReviewUseCase struct {
	repo repository.ReviewRepository
}

func NewReviewUseCase(repository repository.ReviewRepository) *ReviewUseCase {
	return &ReviewUseCase{repo: repository}
}

func (r *ReviewUseCase) AddReview(userID int, productID string, Rating float64, Comment string) (models.ReviewResponse, error) {

	if Rating > 5 || Rating <= 0 {
		return errors.New("Invalid Rating")
	}

	if Comment == "" {
		return errors.New("Please provide a Comment")
	}
	isReviewed, err := r.repo.IsProductReviewedByUser(userID, productID)
	if err != nil {
		return err
	}
	if isReviewed {
		return errors.New("You have already reviewed this product")
	}
	err = r.repo.AddReview(userID, productID, Rating, Comment)
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
	exists, err := r.repo.DoesReviewExist(reviewID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("review not found")
	}
	err = r.repo.DeleteReview(reviewID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReviewUseCase) GetAverageRating(productID string) (float64, error) {
	return r.repo.GetAverageRating(productID)
}
