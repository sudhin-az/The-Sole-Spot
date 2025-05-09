package usecase

import (
	"ecommerce_clean_arch/pkg/repository"
	"ecommerce_clean_arch/pkg/utils/models"
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
		return models.ReviewResponse{}, errors.New("Invalid Rating")
	}

	if Comment == "" {
		return models.ReviewResponse{}, errors.New("Please provide a Comment")
	}
	isReviewed, err := r.repo.IsProductReviewedByUser(userID, productID)
	if err != nil {
		return models.ReviewResponse{}, err
	}
	if isReviewed {
		return models.ReviewResponse{}, errors.New("You have already reviewed this product")
	}
	review, err := r.repo.AddReview(userID, productID, Rating, Comment)
	if err != nil {
		return models.ReviewResponse{}, err
	}
	return models.ReviewResponse{
		ID:        review.ID,
		UserID:    review.UserID,
		ProductID: review.ProductID,
		Rating:    review.Rating,
		Comment:   review.Comment,
	}, nil
}

func (r *ReviewUseCase) GetReviewsByProductID(productID string) ([]models.ReviewResponse, error) {

	exists, err := r.repo.DoesProductExist(productID)
	if err != nil {
		return []models.ReviewResponse{}, err
	}
	if !exists {
		return []models.ReviewResponse{}, errors.New("review not found in this product")
	}
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
	exists, err := r.repo.DoesProductExist(productID)
	if err != nil {
		return 0.0, err
	}
	if !exists {
		return 0.0, errors.New("review not found in this product")
	}
	avgRating, err := r.repo.GetAverageRating(productID)
	if err != nil {
		return 0.0, err
	}
	return avgRating, nil
}
