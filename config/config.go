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
		//set .env variables manually in user console if use heroku
		log.Println(err)
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
