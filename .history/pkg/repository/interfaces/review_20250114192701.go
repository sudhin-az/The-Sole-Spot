package interfaces

type ReviewRepository interface {
	AddReview(userID string, productID string, Rating int, Comment string) error
	GetReviewsByProductID(productID string) models.
}
