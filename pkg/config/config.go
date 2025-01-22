package config

import (
	"fmt"

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

	viper.AddConfigPath("C:/Users/hp/Documents/GO/Ecommerce_clean_architecture/cmd1")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return config, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("errrrr", err)
	}
	return config, nil
}
