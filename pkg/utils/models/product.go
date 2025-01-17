package models

type AddProduct struct {
	ID         int     `json:"id"`
	CategoryID int     `json:"category_id"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Stock      int     `json:"stock"`
	Price      float64 `json:"price"`
}
type ProductResponse struct {
	ID         int     `json:"id" `
	CategoryID int     `json:"category_id"`
	Name       string  `json:"name" `
	Quantity   int     `json:"quantity"`
	Stock      int     `json:"stock"`
	Price      float64 `json:"price"`
}
type ProductDetails struct {
	Name       string  `json:"name"`
	TotalPrice float64 `json:"total_price"`
	Price      float64 `json:"price" `
	Total      float64 `json:"total"`
	Quantity   int     `json:"quantity"`
}
