package main

import (
	"ecommerce_clean_architecture/pkg/api"
	"ecommerce_clean_architecture/pkg/api/handlers"
	"ecommerce_clean_architecture/pkg/config"
	"ecommerce_clean_architecture/pkg/db"
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/usecase"
	"fmt"
	"log"
	"os"

	"github.com/subosito/gotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	// Load environment variables from .env file
	err := gotenv.Load("config.env")
	if err != nil {
		fmt.Println("Vivek!", err)
		log.Fatalf("Error loading .env file")
	}

	// Load configuration
	config, configErr := config.LoadConfig()
	if configErr != nil {
		fmt.Println("Amaan!", configErr)
		log.Fatal("Cannot load config:", configErr)
	}
	fmt.Println("Loaded config:", config)

	// Initialize database connection
	database, dbErr := db.ConnectDatabase(config)
	if dbErr != nil {
		fmt.Println("Niketh!", err)
		log.Fatal("Cannot load database:", dbErr)
	}
	fmt.Println("Database connected:", database)

	// Initialize repositories, use cases, and handlers
	userRepo := repository.NewUserRepository(database)
	userUseCase := usecase.NewUserUseCase(*userRepo)
	userHandler := handlers.NewUserHandler(*userUseCase)

	adminRepo := repository.NewAdminRepository(database)
	adminUseCase := usecase.NewAdminUseCase(*adminRepo)
	adminHandler := handlers.NewAdminHandler(*adminUseCase)

	categoryRepo := repository.NewCategoryRepository(database)
	categoryUseCase := usecase.NewCategoryUseCase(*categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(*categoryUseCase)

	productRepo := repository.NewProductRepository(database)
	productUseCase := usecase.NewProductUseCase(*productRepo)
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
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	authUseCase := usecase.NewAuthUseCase(*userRepo, oauthConfig)
	authHandler := handlers.NewAuthHandler(authUseCase)

	//Initialize the HTTP Server with the user handler
	server := api.NewServerHTTP(userHandler, authHandler, adminHandler, categoryHandler, productHandler, reviewHandler, cartHandler, orderHandler,
		paymentHandler, walletHandler, wishlistHandler, couponHandler)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	//Start the Server on port 8080
	log.Printf("Starting server on port %s...\n", port)
	server.Start(port)
}
