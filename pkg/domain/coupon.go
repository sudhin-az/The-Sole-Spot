package domain

import "time"

type Coupons struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CouponCode      string    `json:"coupon_code" gorm:"unique;not null"`
	Discount        uint      `json:"discount"`
	MinimumRequired uint      `json:"minimum_required"`
	MaximumAllowed  uint      `json:"maximum_allowed"`
	MaximumUsage    uint      `json:"maximum_usage"`
	StartDate       time.Time `json:"createTime,omitempty"`
	EndDate         time.Time `json:"expire_date"`
	ISActive        bool      `json:"is_active,omitempty"`
}
