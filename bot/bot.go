package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"tgbot/config"
	"time"
)

type TermTimer struct {
	timer *time.Timer
	term  chan bool
}

type Bot struct {
	client         *tgbotapi.BotAPI
	chatToScreener map[int64]map[string]*TermTimer
}

func NewBot() *Bot {
	var token = config.GetToken()
	var httpClient = getHttpClient()

	client, err := tgbotapi.NewBotAPIWithClient(token, httpClient)

	if err != nil {
		log.Panic(err)
	}

	chatToScreener := make(map[int64]map[string]*TermTimer)

	log.Printf("Authorized on account %s", client.Self.UserName)

	return &Bot{client, chatToScreener}
}

func (b *Bot) handleUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.client.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
	}
	updates.Clear()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		chatId := update.Message.Chat.ID

		msg := tgbotapi.NewMessage(chatId, "")
		msg.ParseMode = "Markdown"

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "help":
				b.handleHelp(&msg)

			case "status":
				b.handleStatus(&msg)

			case "screen":
				go b.handleInstantScreen(update, &msg)

			case "every":
				b.handleTimeScreen(update, &msg)

			case "clear":
				b.handleClear(&msg)

			default:
				msg.Text = "I don't know that command"
			}
			b.client.Send(msg)
		} else {
			msg.Text = "Please read */help* before you start."
			b.client.Send(msg)
		}
	}
}

func (b *Bot) Start() {
	b.handleUpdates()
}

func (b *Bot) handleHelp(msg *tgbotapi.MessageConfig) {
	var help = []string{
		"Hi, I'm *Sceeny*, I can make some site screenshots for you.",

		"You can use the following *commands*:\n",

		"*Make screen:*",
		"/screen https://google.com\n",

		"*Make screen every minute(m), hour(h), day(d):*",
		"/every 1h https://google.com\n",

		"*Clear all timing screeners:*",
		"/clear\n",
	}
	msg.Text = strings.Join(help, "\n")
}

func (b *Bot) handleInstantScreen(update tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	chatId := msg.ChatID
	textArray := strings.Split(update.Message.Text, " ")
	if len(textArray) != 2 {
		msg.Text = "Bad input string."
		return
	}

	b.sendScreen(chatId, &Screener{
		duration: 0,
		url:      parseUrl(textArray[1]),
		chatId:   chatId,
	})
}

func (b *Bot) handleTimeScreen(update tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	chatId := msg.ChatID
	textArray := strings.Split(update.Message.Text, " ")
	if len(textArray) != 3 {
		msg.Text = "Bad input string."
		return
	}

	screener, err := NewScreener(chatId, textArray[1], textArray[2])
	if err != nil {
		msg.Text = err.Error()
	} else {
		go b.createScreenerJob(screener)
		msg.Text = "Screener set *successfully*!"
	}
}

func (b *Bot) createScreenerJob(screener *Screener) {
	timer := time.NewTimer(screener.duration)
	log.Printf("[create:new] %s", screener)

	if b.chatToScreener[screener.chatId] == nil {
		b.chatToScreener[screener.chatId] = make(map[string]*TermTimer)
	}
	chatMap := b.chatToScreener[screener.chatId]

	shaStr := strconv.Itoa(rand.Int())
	termTimer := &TermTimer{timer, make(chan bool)}
	chatMap[shaStr] = termTimer

	for {
		select {
		case endTime := <-timer.C:
			delete(chatMap, shaStr)
			log.Printf("[make:screen] %s, time: %s", screener, endTime)

			b.sendScreen(screener.chatId, screener)

			timer.Reset(screener.duration)
		case <-termTimer.term:
			log.Printf("[stop:screener] %s", screener)
			return
		}
	}
}

func (b *Bot) sendScreen(chatId int64, screener *Screener) {
	filePath, err := screener.MakeScreen()
	if err != nil {
		log.Println(`[error] make screen error. filepath: `, screener, err)

		msg := tgbotapi.NewMessage(screener.chatId, "Something went wrong with your site. You can try to add *https://www* prefix, some times it's helpful")
		msg.ParseMode = "Markdown"

		b.client.Send(msg)
		return
	}

	defer removeFile(filePath)

	photo := tgbotapi.NewDocumentUpload(chatId, filePath)
	b.client.Send(photo)
}

func (b *Bot) handleStatus(msg *tgbotapi.MessageConfig) {
	msg.Text += fmt.Sprintf("Active timers: *%d*\n", len(b.chatToScreener[msg.ChatID]))
}

func (b *Bot) handleClear(msg *tgbotapi.MessageConfig) {
	screenersMap := b.chatToScreener[msg.ChatID]

	if len(screenersMap) == 0 {
		msg.Text = "There are no active screeners now."
	} else {
		msg.Text = fmt.Sprintf("Stoping %d screeners.", len(screenersMap))
	}
	stopTimers(screenersMap)
}

func stopTimers(timersMaps ...map[string]*TermTimer) {
	for _, timersMap := range timersMaps {
		for key, termTimer := range timersMap {
			termTimer.timer.Stop()
			termTimer.term <- true
			delete(timersMap, key)
		}
	}
}
