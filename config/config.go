package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Token string
	Proxy string
}

var config = &Config{}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.Token = os.Getenv("BOT_TOKEN")
	config.Proxy = os.Getenv("PROXY")

	err = os.Setenv("TZ", "UTC")
	if err != nil {
		log.Println(err)
	}
}

func GetToken() string {
	return config.Token
}

func GetProxy() string {
	return config.Token
}
