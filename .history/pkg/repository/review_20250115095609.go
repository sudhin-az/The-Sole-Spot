package repository

import (
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	DB *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{DB: db}
}

func (r *ReviewRepository) AddReview(userID int, productID uint, Rating float64, Comment string) (models.ReviewResponse, error) {
	query := "INSERT INTO reviews (product_id, user_id, rating, comment) VALUES (?, ?, ?, ?)"
	if err := r.DB.Exec(query, productID, userID, Rating, Comment).Error; err != nil {
		return models.ReviewResponse{}, errors.New("encountered an issue while inserting into reviews")
	}
	return models.ReviewResponse{}, nil
}
func (r *ReviewRepository) IsProductReviewedByUser(userID int, productID uint) (bool, error) {
	var count int64
	err := r.DB.Raw("SELECT COUNT(*) FROM reviews WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *ReviewRepository) GetReviewsByProductID(productID string) ([]models.ReviewResponse, error) {
	var reviews []models.ReviewResponse
	err := r.DB.Raw("SELECT * FROM reviews WHERE product_id = ?", productID).Scan(&reviews).Error
	if err != nil {
		return []models.ReviewResponse{}, err
	}
	return reviews, nil
}
func (r *ReviewRepository) DoesReviewExist(reviewID string) (bool, error) {
	var count int64
	err := r.DB.Raw("SELECT COUNT(*) FROM reviews WHERE id = ?", reviewID).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *ReviewRepository) DeleteReview(reviewID string) error {
	err := r.DB.Raw("DELETE FROM reviews WHERE id = ?", reviewID).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ReviewRepository) GetAverageRating(productID string) (float64, error) {
	var AvgRating float64
	err := r.DB.Raw("SELECT AVG(rating) FROM reviews WHERE product_id = ?", productID).Scan(&AvgRating).Error
	if err != nil {
		return 0, err
	}
	return AvgRating, nil
}
