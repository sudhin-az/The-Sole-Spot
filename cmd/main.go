package main

import (
	"ecommerce_clean_arch/pkg/config"
	"ecommerce_clean_arch/pkg/db"
	"ecommerce_clean_arch/pkg/di"
	"log"
	"os"

	"github.com/subosito/gotenv"
)

// @title The-Sole-Spot API
// @description This is the API documentation for The-Sole-Spot application.
// @version 1.0
// @host localhost:8080
// @BasePath /

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
	log.Println("Loaded config:", config)

	// Initialize database connection
	database, dbErr := db.ConnectDatabase(config)
	if dbErr != nil {
		log.Fatal("Cannot load database:", dbErr)
	}
	log.Println("Database connected:", database)

	server, err := di.InitializeAPI(config)
	if err != nil {
		log.Fatal("API initialization failed:", err)
	}
	log.Println("Dependencies initialized successfully")

	server, err = di.InitializeServer(config)
	if err != nil {
		log.Println("Server initialization failed")
	}
	log.Println("Server Dependencies initialized successfully")

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	//Start the Server on port 8080
	log.Printf("Starting server on port %s...\n", port)
	server.Start(port)
}
