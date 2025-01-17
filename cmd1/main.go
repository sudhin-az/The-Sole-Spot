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
		log.Fatalf("Error loading .env file")
	}

	// Load configuration
	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("Cannot load config:", configErr)
	}
	fmt.Println("Loaded config:", config)

	// Initialize database connection
	database, dbErr := db.ConnectDatabase(config)
	if dbErr != nil {
		log.Fatal("Cannot load database:", dbErr)
	}
	fmt.Println("Database connected:", database)

	// Initialize repositories, use cases, and handlers
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
	productUseCase := usecase.NewProductUseCase(*productRepo)
	productHandler := handlers.NewProductHandler(*productUseCase)

	cartRepo := repository.NewCartRepository(database)
	productRepo = repository.NewProductRepository(database)
	cartUseCase := usecase.NewCartUseCase(*cartRepo, *productRepo)
	cartHandler := handlers.NewCartHandler(*cartUseCase)

	orderRepo := repository.NewOrderRepository(database)
	orderUseCase := usecase.NewOrderUseCase(*orderRepo)
	orderHandler := handlers.NewOrderHandler(*orderUseCase)

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
	server := api.NewServerHTTP(userHandler, authHandler, adminHandler, categoryHandler, productHandler, cartHandler, orderHandler)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	//Start the Server on port 8080
	log.Printf("Starting server on port %s...\n", port)
	server.Start(port)
}
