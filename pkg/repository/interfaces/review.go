package interfaces

import "ecommerce_clean_arch/pkg/utils/models"

type ReviewRepository interface {
	AddReview(userID int, productID string, Rating float64, Comment string) (models.ReviewResponse, error)
	GetReviewsByProductID(productID string) ([]models.ReviewResponse, error)
	IsProductReviewedByUser(userID string, productID string) (bool, error)
	DoesProductExist(productID string) (bool, error)
	DoesReviewExist(reviewID string) (bool, error)
	DeleteReview(reviewID string) error
	GetAverageRating(productID string) (float64, error)
}
