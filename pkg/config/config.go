package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost                      string `mapstructure:"DB_HOST" validate:"required"`
	DBName                      string `mapstructure:"DB_NAME" validate:"required"`
	DBUser                      string `mapstructure:"DB_USER" validate:"required"`
	DBPort                      string `mapstructure:"DB_PORT" validate:"required"`
	DBPassword                  string `mapstructure:"DB_PASSWORD" validate:"required"`
	ClientID                    string `mapstructure:"GOOGLE_CLIENT_ID" validate:"required"`
	ClientSecret                string `mapstructure:"GOOGLE_CLIENT_SECRET" validate:"required"`
	RedirectURL                 string `mapstructure:"GOOGLE_REDIRECT_URL" validate:"required"`
	KEY_ID_FOR_RAYZORPAY        string `mapstructure:"KEY_ID_FOR_RAYZORPAY"`
	SECRET_KEY_ID_FOR_RAYZORPAY string `mapstructure:"SECRET_KEY_ID_FOR_RAYZORPAY"`
}

func LoadConfig() (Config, error) {
	var config Config

	viper.AddConfigPath("C:/Users/hp/Desktop/Ecommerce_clean_architecture/cmd1")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
		return config, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Println("errrrr", err)
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
