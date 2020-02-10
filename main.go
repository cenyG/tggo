package main

import (
	"tgbot.go/bot"
	_ "tgbot.go/config"
)

func main() {
	myBot := bot.NewBot()
	myBot.Start()
}
