package domain

type Coupon struct {
	ID                 int    `json:"id" gorm:"uniquekey; not null"`
	CouponCode         string `json:"coupon_code" gorm:"coupon_code"`
	DiscountPercentage int    `json:"discount_percentage"`
	MinumumPrice       int    `json:"minumum_price"`
	Description        string `json:"description"`
	IsActive           bool
}
type UsedCoupon struct {
	ID       int    `json:"id" gorm:"uniquekey not null"`
	CouponID int    `json:"coupon_id"`
	Coupon   Coupon `json:"-" gorm:"foreignkey:CouponID"`
	UserID   int    `json:"user_id"`
	Users    Users  `json:"-" gorm:"foreignkey:UserID"`
	IsActive bool   `json:"used"`
}
