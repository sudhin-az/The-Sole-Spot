package main

import (
	"ecommerce_clean_arch/pkg/config"
	"ecommerce_clean_arch/pkg/db"
	"ecommerce_clean_arch/pkg/di"
	"fmt"
	"log"
	"os"
)

// @title The-Sole-Spot API
// @description This is the API documentation for The-Sole-Spot application.
// @version 1.0
// @host {{.ServerHost}}
// @BasePath /

func main() {
	// Load configuration
	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("Cannot load config:", configErr)
	}
	log.Println("Config loaded successfully")

	// Initialize database connection
	database, dbErr := db.ConnectDatabase(config)
	if dbErr != nil {
		log.Fatal("Cannot connect to database:", dbErr)
	}
	log.Println("Database connected successfully")
	fmt.Println(database)

	// Initialize API dependencies
	server, err := di.InitializeAPI(config)
	if err != nil {
		log.Fatal("API initialization failed:", err)
	}
	log.Println("API Dependencies initialized successfully")

	// Initialize Server dependencies
	server, err = di.InitializeServer(config)
	if err != nil {
		log.Fatal("Server initialization failed:", err)
	}
	log.Println("Server Dependencies initialized successfully")

	// Get the port from environment or default to 8080
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Allow external connections in production
	host := "0.0.0.0"
	if os.Getenv("ENV") == "dev" {
		host = "localhost"
	}

	// Start the Server
	address := host + ":" + port
	log.Printf("Starting server on %s...\n", address)
	server.Start(address)
}
