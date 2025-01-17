package interfaces

type OrderRepository interface {
	DoesCartExist(userID int) (bool, error)
}
