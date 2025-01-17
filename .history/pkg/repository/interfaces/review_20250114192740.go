package interfaces

import "ecommerce_clean_architecture/pkg/utils/models"

type ReviewRepository interface {
	AddReview(userID string, productID string, Rating int, Comment string) error
	GetReviewsByProductID(productID string) ([]models.ReviewResponse, error)
}
