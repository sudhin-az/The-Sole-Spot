package models

type ReviewRequest struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment"`
}
type ReviewResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}
