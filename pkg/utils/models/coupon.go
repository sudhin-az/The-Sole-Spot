package models

import "time"

type Coupon struct {
	ID              uint      `json:"coupon_id"`
	CouponCode      string    `json:"coupon_code" validate:"required,alphanum,min=5,max=20"`
	Discount        uint      `json:"discount" validate:"required,min=1,max=100"`
	MinimumRequired uint      `json:"minimum_required" validate:"required,min=0"`
	MaximumAllowed  uint      `json:"maximum_allowed" validate:"required,gtcsfield=MinimumRequired"`
	MaximumUsage    uint      `json:"maximum_usage" validate:"required,min=2"`
	ExpireDateStr   string    `json:"expire_date" validate:"required"`
	ExpireDate      time.Time `json:"_"`
	ISActive        bool      `json:"is_active"`
}

type CouponResponse struct {
	ID              uint      `json:"couponID"`
	CouponCode      string    `json:"coupon_code"`
	Discount        uint      `json:"discount"`
	MinimumRequired uint      `json:"minimum_required"`
	MaximumAllowed  uint      `json:"maximum_allowed"`
	MaximumUsage    uint      `json:"maximum_usage"`
	StartDate       time.Time `json:"createTime,omitempty"`
	EndDate         time.Time `json:"expire_date"`
	ISActive        bool      `json:"is_active,omitempty"`
}
