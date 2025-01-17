package domain

type Review struct {
	ID        int
	User      Users
	UserID    int
	Product   Products
	ProductID int
	Rating    int
	Comment   string
}
