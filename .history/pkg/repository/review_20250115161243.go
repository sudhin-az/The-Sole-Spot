package repository

import (
	"ecommerce_clean_architecture/pkg/domain"
	"ecommerce_clean_architecture/pkg/utils/models"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	DB *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{DB: db}
}

func (r *ReviewRepository) AddReview(userID int, productID string, Rating float64, Comment string) (domain.Review, error) {

	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		return domain.Review{}, fmt.Errorf("invalid productID: %v", err)
	}
	Review := domain.Review{
		UserID:    uint(userID),
		ProductID: uint(productIDUint),
		Rating:    Rating,
		Comment:   Comment,
	}
	err = r.DB.Create(&Review).Error
	if err != nil {
		return domain.Review{}, fmt.Errorf("database error creating review: %w", err)
	}
	return Review, nil
}
func (r *ReviewRepository) IsProductReviewedByUser(userID int, productID string) (bool, error) {
	var count int64
	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		return false, fmt.Errorf("invalid product ID: %w", err)
	}
	err = r.DB.Model(&domain.Review{}).Where("user_id = ? AND product_id = ?", userID, productIDUint).
		Count(&count).Error
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
