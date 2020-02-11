package main

import (
	"tgbot/bot"
	_ "tgbot/config"
)

func main() {
	myBot := bot.NewBot()
	myBot.Start()
}
