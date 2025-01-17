package interfaces

import "ecommerce_clean_architecture/pkg/utils/models"

type ReviewUseCase interface {
	AddReview(userID string, productID string, Rating int, Comment string) error
	GetReviewsByProductID(productID string) ([]models.ReviewResponse, error)
	DeleteReview(reviewID string) error
	GetAverageRating(productID string) (float64, error)
}
