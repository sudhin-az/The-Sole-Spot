package repository

import (
	"ecommerce_clean_arch/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type CouponRepository struct {
	DB *gorm.DB
}

func NewCouponRepository(db *gorm.DB) *CouponRepository {
	return &CouponRepository{
		DB: db,
	}
}

func (coup *CouponRepository) AddCoupon(coupon models.Coupon) (models.CouponResponse, error) {
	var createdCoupon models.CouponResponse

	couponExpireTime := coupon.ExpireDate

	query := `INSERT INTO coupons (coupon_code, discount, minimum_required, maximum_allowed,maximum_usage, start_date, end_date, is_active) 
	VALUES(?, ?, ?, ?, ?, now(), ?, true) RETURNING *`
	result := coup.DB.Raw(query, coupon.CouponCode, coupon.Discount, coupon.MinimumRequired, coupon.MaximumAllowed, coupon.MaximumUsage, couponExpireTime).Scan(&createdCoupon)
	if result.Error != nil {
		return models.CouponResponse{}, errors.New("encountered an issue while creating a new coupon")
	}
	if result.RowsAffected == 0 {
		return models.CouponResponse{}, errors.New("No rows affected")
	}
	return createdCoupon, nil
}

func (coup *CouponRepository) MakeCouponInvalid(id int) error {
	if err := coup.DB.Exec("update coupons set is_active=false where id=$1", id).Error; err != nil {
		return err
	}

	return nil
}

func (coup *CouponRepository) GetAllCoupons() ([]models.CouponResponse, error) {
	var coupons []models.CouponResponse
	err := coup.DB.Raw("select * from coupons").Scan(&coupons).Error
	if err != nil {
		return []models.CouponResponse{}, err
	}
	return coupons, nil
}

func (coup *CouponRepository) CheckCouponExpired(tx *gorm.DB, couponCode string) (models.CouponResponse, error) {
	var couponData models.CouponResponse

	query := "SELECT * FROM coupons WHERE coupon_code = ? AND is_active = 'true'"
	result := coup.DB.Raw(query, couponCode).Scan(&couponData)
	if result.Error != nil {
		return models.CouponResponse{}, errors.New("face some issue while check coupon exist")
	}
	if result.RowsAffected == 0 {
		return models.CouponResponse{}, errors.New("not a valid coupon, better luck next time")
	}
	return couponData, nil
}

func (coup *CouponRepository) GetCouponUsageCount(couponCode string, userID uint) (int, error) {
	var usageCount int
	query := `SELECT COUNT(*) FROM coupons WHERE coupon_code = ? AND user_id = ?`
	err := coup.DB.Raw(query, couponCode, userID).Scan(&usageCount).Error
	if err != nil {
		return 0, errors.New("Could not check coupon usage")
	}
	return usageCount, nil
}

func (coup *CouponRepository) UpdateCouponStatus(couponID, active string) (models.CouponResponse, error) {
	var coupon models.CouponResponse
	var result *gorm.DB

	query := "UPDATE coupons SET is_active = ? WHERE id = ? RETURNING*"
	if active != "" {
		result = coup.DB.Raw(query, active, couponID)
	}
	if result.Error != nil {
		return models.CouponResponse{}, errors.New("face some issue while update coupons status")
	}
	if result.RowsAffected == 0 {
		return models.CouponResponse{}, errors.New("No rows affected")
	}
	return coupon, nil
}
