package interfaces

import "ecommerce_clean_arch/pkg/utils/models"

type ReviewUseCase interface {
	AddReview(userID int, productID string, Rating float64, Comment string) (models.ReviewResponse, error)
	GetReviewsByProductID(productID string) ([]models.ReviewResponse, error)
	DeleteReview(reviewID string) error
	GetAverageRating(productID string) (float64, error)
}
