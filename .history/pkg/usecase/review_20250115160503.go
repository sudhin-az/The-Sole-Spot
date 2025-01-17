package usecase

import (
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
	"fmt"
)

type ReviewUseCase struct {
	repo repository.ReviewRepository
}

func NewReviewUseCase(repository repository.ReviewRepository) *ReviewUseCase {
	return &ReviewUseCase{repo: repository}
}

func (r *ReviewUseCase) AddReview(userID int, productID string, rating float64, comment string) (models.ReviewResponse, error) {
	// Validate rating
	if rating > 5 || rating <= 0 {
		return models.ReviewResponse{}, errors.New("Invalid Rating")
	}

	// Validate comment
	if comment == "" {
		return models.ReviewResponse{}, errors.New("Please provide a Comment")
	}

	// Check if product is already reviewed by user
	isReviewed, err := r.repo.IsProductReviewedByUser(userID, productID)
	if err != nil {
		return models.ReviewResponse{}, fmt.Errorf("error checking review status: %w", err)
	}
	if isReviewed {
		return models.ReviewResponse{}, errors.New("You have already reviewed this product")
	}

	// Add the review
	review, err := r.repo.AddReview(userID, productID, rating, comment)
	if err != nil {
		return models.ReviewResponse{}, fmt.Errorf("error adding review: %w", err)
	}

	// Return success response
	return models.ReviewResponse{
		ID:        review.ID,
		UserID:    review.UserID,
		ProductID: review.ProductID,
		Rating:    review.Rating,
		Comment:   review.Comment,
	}, nil
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
