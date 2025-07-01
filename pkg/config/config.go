package config

import (
	"fmt"
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

	// Get current and parent directories
	currentDir, err := os.Getwd()
	if err != nil {
		return config, fmt.Errorf("failed to get working directory: %v", err)
	}
	parentDir := filepath.Dir(currentDir)

	searchPaths := []string{currentDir, parentDir}

	// Check if env variables are already set
	requiredVars := []string{
		"DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD",
		"GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET", "GOOGLE_REDIRECT_URL",
	}

	allSet := true
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			allSet = false
			fmt.Printf("⚠️ missing env var: %s\n", v)
		}
	}

	if allSet {
		fmt.Println("✅ all required environment variables are already set.")

		// prioritize K8 > DOCKER > localhost
		var host string
		switch {
		case os.Getenv("K8") == "YES":
			host = os.Getenv("DB_HOST")
		case os.Getenv("DOCKER") == "YES":
			host = "host.docker.internal"
		default:
			host = "localhost"
		}

		config = Config{
			DBHost:       host,
			DBName:       os.Getenv("DB_NAME"),
			DBUser:       os.Getenv("DB_USER"),
			DBPort:       os.Getenv("DB_PORT"),
			DBPassword:   os.Getenv("DB_PASSWORD"),
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		}

		fmt.Println(config)
		return config, nil
	}

	// fallback: try loading from .env file
	for _, path := range searchPaths {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(path)
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("📄 loaded config from:", viper.ConfigFileUsed())
			if err := viper.Unmarshal(&config); err != nil {
				return config, fmt.Errorf("failed to unmarshal config: %v", err)
			}
			return config, nil
		}
	}

	return config, fmt.Errorf(".env file not found in current or parent directory")
}
