package usecase

import (
	"ecommerce_clean_arch/pkg/repository"
	"ecommerce_clean_arch/pkg/utils/models"
	"errors"
)

type CouponUseCase struct {
	repository repository.CouponRepository
}

func NewCouponUseCase(repo repository.CouponRepository) *CouponUseCase {
	return &CouponUseCase{
		repository: repo,
	}
}

func (coup *CouponUseCase) AddCoupon(newcoupon models.Coupon) (models.CouponResponse, error) {
	if newcoupon.Discount >= 100 {
		return models.CouponResponse{}, errors.New("Discount percentage cannot be greater than 100")
	}
	coupon, err := coup.repository.AddCoupon(newcoupon)
	if err != nil {
		return models.CouponResponse{}, err
	}
	return coupon, nil
}

func (coup *CouponUseCase) MakeCouponInvalid(id int) error {
	err := coup.repository.MakeCouponInvalid(id)
	if err != nil {
		return err
	}
	return nil
}

func (coup *CouponUseCase) GetAllCoupons() ([]models.CouponResponse, error) {
	coupon, err := coup.repository.GetAllCoupons()
	if err != nil {
		return []models.CouponResponse{}, err
	}
	return coupon, nil
}

func (coup *CouponUseCase) UpdateCouponStatus(couponID, status string) (models.CouponResponse, error) {
	var coupon models.CouponResponse
	var err error
	if status == "active" {
		coupon, err = coup.repository.UpdateCouponStatus(couponID, status)
		if err != nil {
			return models.CouponResponse{}, err
		}
	}
	return coupon, nil
}
