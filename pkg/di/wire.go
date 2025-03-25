package di

import (
	"ecommerce_clean_arch/pkg/api"
	"ecommerce_clean_arch/pkg/api/handlers"
	"ecommerce_clean_arch/pkg/config"
	"ecommerce_clean_arch/pkg/db"
	"ecommerce_clean_arch/pkg/repository"
	"ecommerce_clean_arch/pkg/usecase"
	"log"

	"github.com/google/wire"
)

// Provide all dependencies
func InitializeAPI(cfg config.Config) (*api.ServerHTTP, error) {
	wire.Build(
		db.ConnectDatabase,

		// Repositories
		repository.NewUserRepository,
		repository.NewAdminRepository,
		repository.NewCategoryRepository,
		repository.NewProductRepository,
		repository.NewCartRepository,
		repository.NewWalletRepository,
		repository.NewCouponRepository,
		repository.NewWishlistRepository,
		repository.NewOrderRepository,
		repository.NewReviewRepository,
		repository.NewPaymentRepository,

		// Use Cases
		usecase.NewUserUseCase,
		usecase.NewAdminUseCase,
		usecase.NewCategoryUseCase,
		usecase.NewProductUseCase,
		usecase.NewCartUseCase,
		usecase.NewWalletUseCase,
		usecase.NewCouponUseCase,
		usecase.NewWishlistUseCase,
		usecase.NewOrderUseCase,
		usecase.NewReviewUseCase,
		usecase.NewPaymentUsecase,

		// Handlers
		handlers.NewUserHandler,
		handlers.NewAuthHandler,
		handlers.NewAdminHandler,
		handlers.NewCategoryHandler,
		handlers.NewProductHandler,
		handlers.NewCartHandler,
		handlers.NewWalletHandler,
		handlers.NewCouponHandler,
		handlers.NewWishlistHandler,
		handlers.NewOrderHandler,
		handlers.NewReviewHandler,
		handlers.NewPaymentHandler,

		// Server
		api.NewServerHTTP,
	)
	log.Println("Dependencies initialized successfully")
	return &api.ServerHTTP{}, nil
}
