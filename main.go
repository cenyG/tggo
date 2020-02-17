package main

import (
	"os"
	"tgbot/bot"
	_ "tgbot/config"
	"tgbot/heroku"
)

func main() {
	if os.Getenv(`HEROKU`) == `true` {
		go heroku.MockServer()
	}

	myBot := bot.NewBot()
	myBot.Start()
}
