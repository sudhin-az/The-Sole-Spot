package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost       string `mapstructure:"DB_HOST" validate:"required"`
	DBName       string `mapstructure:"DB_NAME" validate:"required"`
	DBUser       string `mapstructure:"DB_USER" validate:"required"`
	DBPort       string `mapstructure:"DB_PORT" validate:"required"`
	DBPassword   string `mapstructure:"DB_PASSWORD" validate:"required"`
	ClientID     string `mapstructure:"GOOGLE_CLIENT_ID" validate:"required"`
	ClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET" validate:"required"`
	RedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL" validate:"required"`
}

func LoadConfig() (Config, error) {
	var config Config

	// Determine the environment (dev or prod)
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // Default to dev
	}
	log.Printf("Running in %s mode", env)

	// Get the config path based on the environment
	var configPath string
	if env == "production" {
		configPath = os.Getenv("CONFIG_PATH_PROD")
	} else {
		configPath = os.Getenv("CONFIG_PATH_DEV")
	}

	if configPath == "" {
		log.Fatalf("CONFIG_PATH not set for environment: %s", env)
	}

	// Log the raw config path before absolute conversion
	log.Printf("Raw config path: %s", configPath)

	// Ensure that the path is correctly formatted
	configPath = filepath.Clean(configPath)

	log.Printf("Cleaned config path: %s", configPath)

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("error unmarshalling config: %v", err)
	}

	return config, nil
}

func LoadRazorpayConfig() (string, string) {
	key := os.Getenv("KEY_ID_FOR_RAYZORPAY")
	secret := os.Getenv("SECRET_KEY_ID_FOR_RAYZORPAY")
	if key == "" || secret == "" {
		log.Fatal("Razorpay credentials are not set in the environment")
	}
	return key, secret
}
