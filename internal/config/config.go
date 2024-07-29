package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Configurations for the application
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	BaseURL             string        `mapstructure:"BASE_URL"`
	GrantType           string        `mapstructure:"GRANT_TYPE"`
	Scope               string        `mapstructure:"SCOPE"`
	Login               string        `mapstructure:"CLIENT_ID"`
	Password            string        `mapstructure:"CLIENT_SECRET"`
	ShopID              string        `mapstructure:"SHOP_ID"`
	TerminalID          string        `mapstructure:"TERMINAL_ID"`
	EPAYURL             string        `mapstructure:"EPAY_URL"`
	EPAYLogin           string        `mapstructure:"EPAY_LOGIN"`
	EPAYPassword        string        `mapstructure:"EPAY_PASSWORD"`
	EPAYOAuthURL        string        `mapstructure:"EPAY_OAUTH_URL"`
	EPAYPaymentPageURL  string        `mapstructure:"EPAY_PAYMENT_PAGE_URL"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
