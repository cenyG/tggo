package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Token string
}

var config = &Config{}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.Token = os.Getenv("BOT_TOKEN")
}

func GetToken() string {
	return config.Token
}
