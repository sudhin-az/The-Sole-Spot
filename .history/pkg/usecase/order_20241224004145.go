package usecase

type OrderUseCase struct {
	orderRepository interfaces.OrderRepository
	userRepository  interfaces.UserRepository
	cartRepository  interfaces.CartRepository
}
