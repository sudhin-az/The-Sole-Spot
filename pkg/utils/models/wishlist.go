package models

type WishlistRequest struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
}

type WishListResponse struct {
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
}
