package interfaces

import "ecommerce_clean_arch/pkg/utils/models"

type CouponUseCase interface {
	AddCoupon(coupon models.Coupon) (models.CouponResponse, error)
	MakeCouponInvalid(id int) error
	GetAllCoupons() ([]models.CouponResponse, error)
	UpdateCouponStatus(couponID, active string) (models.CouponResponse, error)
}
