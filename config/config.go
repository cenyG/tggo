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
		log.Println("error loading .env file", err)
	}

	config.Token = os.Getenv("BOT_TOKEN")
	config.Proxy = os.Getenv("PROXY")
}

func GetToken() string {
	return config.Token
}

func GetProxy() string {
	return config.Proxy
}
