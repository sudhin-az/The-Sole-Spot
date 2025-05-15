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

    // Get the current working directory
    workDir, err := os.Getwd()
    if err != nil {
        return config, fmt.Errorf("error getting working directory: %v", err)
    }

    // Move up one directory to reach the project root
    projectRoot := filepath.Dir(workDir)

    // Setup Viper
    viper.SetConfigName(".env")        // name of config file (without extension)
    viper.SetConfigType("env")         // type of config file
    viper.AddConfigPath(projectRoot)    // path to look for the config file in
    viper.AddConfigPath(".")           // optionally look for config in the working directory
    viper.AutomaticEnv()               // read in environment variables that match

    // Read the config file
    if err := viper.ReadInConfig(); err != nil {
        return config, fmt.Errorf("error reading config file: %v", err)
    }

    // Unmarshal the config
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
