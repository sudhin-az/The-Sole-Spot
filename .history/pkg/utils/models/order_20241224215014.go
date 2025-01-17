package models

type OrderFromCart struct {
	PaymentID uint `json:"payment_id" binding:"required"`
	AddressID uint `json:"address_id" binding:"required"`
}

type OrderIncoming s
