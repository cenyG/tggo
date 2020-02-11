package bot

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"tgbot/config"
	"time"
)

type ConcurrentSlice struct {
	sync.RWMutex
	items []interface{}
}

type Bot struct {
	client             *tgbotapi.BotAPI
	chatToTimezone     map[int64]int64
	chatToTimer        map[int64]map[string]*time.Timer
	chatToNotification map[int64]map[string]*time.Timer
}

func NewBot() *Bot {
	var token = config.GetToken()
	var httpClient = getHttpClient()

	client, err := tgbotapi.NewBotAPIWithClient(token, httpClient)
	//client.Debug = true
	if err != nil {
		log.Panic(err)
	}

	chatToTimezone := make(map[int64]int64)
	chatToTimer := make(map[int64]map[string]*time.Timer)
	chatToNotification := make(map[int64]map[string]*time.Timer)

	log.Printf("Authorized on account %s", client.Self.UserName)

	return &Bot{client, chatToTimezone, chatToTimer, chatToNotification}
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

			case "timezone":
				b.handleTimezone(update, &msg)

			case "set":
				b.handleSet(update, &msg)

			case "timer":
				b.handleTimer(update, &msg)

			case "clear":
				b.handleClear(update, &msg)

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
		"This BOT will make notifications for you.",
		"You can use the following *commands*:\n",
		"_Simple notifications for today:_",
		"*/set HH:mm text*\n",
		"_Every day notifications:_",
		"*/set every HH:mm text*\n",
		"_Exact day notifications:_",
		"*/set DD/MM/YYYY HH:mm text*\n",
		"_Timer:_",
		"*/timer mm:ss text*\n",
		"*DD/MM/YYYY* - `(day/month/year)` format.",
		"For example: *02/11/2020*",
		"*HH:mm:ss* - `(hour/minute/seconds)` format.",
		"For example: *22:59*\n",
		"* ⚠️⚠️⚠️ Don't forget to set your UTC /timezone*",
		"For example: */timezone +3*",
	}
	msg.Text = strings.Join(help, "\n")
}

func (b *Bot) handleTimezone(update tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	res := strings.Split(update.Message.Text, " ")
	if len(res) != 2 {
		msg.Text = "*Wrong timezone format*. Please use format like : */timezone -4*"
		return
	}
	tz, err := timezoneParse(res[1])
	if err != nil {
		msg.Text = "*Wrong timezone*. Please use format like : */timezone -4*"
		return
	}
	b.chatToTimezone[msg.ChatID] = tz
	msg.Text = "Timezone set *success*!"
}

func (b *Bot) handleSet(update tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	chatId := msg.ChatID
	textArray := strings.Split(update.Message.Text, " ")

	utcDiff, _ := b.chatToTimezone[chatId]

	notification, err := NewNotification(textArray[1:], chatId, utcDiff)
	if err != nil {
		msg.Text = err.Error()
	} else {
		go b.createNotificationJob(notification)
		msg.Text = "Notification set *successfully*!"
	}
}

func (b *Bot) handleTimer(update tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	chatId := msg.ChatID
	textArray := strings.Split(update.Message.Text, " ")

	myTimer, err := NewTimer(textArray[1:], chatId)
	if err != nil {
		msg.Text = err.Error()
	} else {
		go b.createTimerJob(myTimer)
		msg.Text = "Timer set *successfully*!"
	}
}

func (b *Bot) createTimerJob(mTimer *Timer) {
	timer := time.NewTimer(mTimer.duration)
	log.Printf("[Set] %s", mTimer)

	if b.chatToTimer[mTimer.chatId] == nil {
		b.chatToTimer[mTimer.chatId] = make(map[string]*time.Timer)
	}
	chatMap := b.chatToTimer[mTimer.chatId]

	shaStr := shaString(strconv.Itoa(rand.Int()))
	chatMap[shaStr] = timer

	timerEnd := <-timer.C
	delete(chatMap, shaStr)

	log.Printf("[End] %s, time: %s", mTimer, timerEnd)

	msg := tgbotapi.NewMessage(mTimer.chatId, mTimer.text)
	b.client.Send(msg)

	if mTimer.repeat {
		go b.createTimerJob(mTimer)
	}
}

func (b *Bot) createNotificationJob(notification *Notification) {
	timer := time.NewTimer(notification.duration)
	log.Printf("[Set] %s", notification)

	if b.chatToNotification[notification.chatId] == nil {
		b.chatToNotification[notification.chatId] = make(map[string]*time.Timer)
	}
	chatMap := b.chatToNotification[notification.chatId]

	shaStr := shaString(strconv.Itoa(rand.Int()))
	chatMap[shaStr] = timer

	timerEnd := <-timer.C
	delete(chatMap, shaStr)

	log.Printf("[End] %s, time: %s", notification, timerEnd)

	msg := tgbotapi.NewMessage(notification.chatId, notification.text)
	b.client.Send(msg)

	if notification.repeat {
		notification.duration = 24 * time.Hour
		go b.createNotificationJob(notification)
	}
}

func shaString(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)

	return hex.EncodeToString(bs)
}

func (b *Bot) handleStatus(msg *tgbotapi.MessageConfig) {
	msg.Text = fmt.Sprintf("Timezone: *UTC%+d*\n", b.chatToTimezone[msg.ChatID])
	msg.Text += fmt.Sprintf("Active timers: *%d*\n", len(b.chatToTimer[msg.ChatID]))
	msg.Text += fmt.Sprintf("Active Notifications: *%d*\n", len(b.chatToNotification[msg.ChatID]))
}

func (b *Bot) handleClear(update tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	timersMap := b.chatToTimer[msg.ChatID]
	for key, timer := range timersMap {
		timer.Stop()
		delete(timersMap, key)
	}

	notifiersMap := b.chatToNotification[msg.ChatID]
	for key, timer := range notifiersMap {
		timer.Stop()
		delete(notifiersMap, key)
	}
}
