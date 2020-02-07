package main

import (
	"tgbot.go/bot"
)

func main() {
	myBot := bot.NewBot()
	myBot.Start()
}

func getNotification(text string, chatId int64) (*bot.Notification, error) {
	return bot.NewNotification(text, chatId)
}

func getTimer(text string, chatId int64) (*bot.Timer, error) {
	return bot.NewTimer(text, chatId)
}
