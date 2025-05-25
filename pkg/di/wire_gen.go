package di

import (
	"ecommerce_clean_arch/pkg/api"
	"ecommerce_clean_arch/pkg/api/handlers"
	"ecommerce_clean_arch/pkg/config"
	"ecommerce_clean_arch/pkg/db"
	"ecommerce_clean_arch/pkg/repository"
	"ecommerce_clean_arch/pkg/usecase"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func InitializeServer(cfg config.Config) (*api.ServerHTTP, error) {
	database, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}

	// Initialization
	userRepo := repository.NewUserRepository(database)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handlers.NewUserHandler(*userUseCase)

	adminRepo := repository.NewAdminRepository(database)
	adminUseCase := usecase.NewAdminUseCase(*adminRepo)
	adminHandler := handlers.NewAdminHandler(*adminUseCase)

	categoryRepo := repository.NewCategoryRepository(database)
	categoryUseCase := usecase.NewCategoryUseCase(*categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(*categoryUseCase)

	productRepo := repository.NewProductRepository(database)
	productUseCase := usecase.NewProductUseCase(productRepo)
	productHandler := handlers.NewProductHandler(*productUseCase)

	cartRepo := repository.NewCartRepository(database)
	cartUseCase := usecase.NewCartUseCase(*cartRepo, *productRepo, *categoryRepo)
	cartHandler := handlers.NewCartHandler(*cartUseCase)

	walletRepo := repository.NewWalletRepository(database)
	walletUseCase := usecase.NewWalletUseCase(*walletRepo)
	walletHandler := handlers.NewWalletHandler(*walletUseCase)

	couponRepo := repository.NewCouponRepository(database)
	couponUseCase := usecase.NewCouponUseCase(*couponRepo)
	couponHandler := handlers.NewCouponHandler(*couponUseCase)

	wishlistRepo := repository.NewWishlistRepository(database)
	wishlistUseCase := usecase.NewWishlistUseCase(*wishlistRepo)
	wishlistHandler := handlers.NewWishlistHandler(*wishlistUseCase)

	orderRepo := repository.NewOrderRepository(database)
	orderUseCase := usecase.NewOrderUseCase(*orderRepo, *userRepo, *cartRepo, *walletRepo, *walletUseCase, *couponRepo)
	orderHandler := handlers.NewOrderHandler(*orderUseCase)

	reviewRepo := repository.NewReviewRepository(database)
	reviewUseCase := usecase.NewReviewUseCase(*reviewRepo)
	reviewHandler := handlers.NewReviewHandler(*reviewUseCase)

	paymentRepo := repository.NewPaymentRepository(database)
	paymentUseCase := usecase.NewPaymentUsecase(*paymentRepo)
	paymentHandler := handlers.NewPaymentHandler(*paymentUseCase)

	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	authUseCase := usecase.NewAuthUseCase(*userRepo, oauthConfig)
	authHandler := handlers.NewAuthHandler(authUseCase)

	server := api.NewServerHTTP(userHandler, authHandler, adminHandler, categoryHandler, productHandler, reviewHandler, cartHandler, orderHandler,
		paymentHandler, walletHandler, wishlistHandler, couponHandler)

	return server, nil
}
