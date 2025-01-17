package repository

import (
	"ecommerce_clean_architecture/pkg/utils/models"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	DB *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{DB: db}
}

func (r *ReviewRepository) AddReview(userID string, productID string, Rating int, Comment string) error {
	var review models.ReviewResponse
	err := r.DB.Raw("INSERT INTO reviews(user_id, product_id, rating, comment) VALUES(?,?,?,?)",
		review.UserID, review.ProductID, review.Rating, review.Comment).Scan(&review).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ReviewRepository) GetReviewsByProductID(productID string) ([]models.ReviewResponse, error) {
	var reviews []models.ReviewResponse
	err := r.DB.Raw("SELECT * FROM reviews WHERE product_id = ?", productID).Scan(&reviews).Error
	if err != nil {
		return []models.ReviewResponse{}, err
	}
	return reviews, nil
}
