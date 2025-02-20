package interfaces

import "ecommerce_clean_architecture/pkg/utils/models"

type CouponRepository interface {
	AddCoupon(coupon models.Coupon) (models.CouponResponse, error)
	MakeCouponInvalid(id int) error
	GetAllCoupons() ([]models.CouponResponse, error)
	CheckCouponExpired(couponCode string) (models.CouponResponse, error)
	UpdateCouponStatus(couponID, active string) (models.CouponResponse, error)
}
