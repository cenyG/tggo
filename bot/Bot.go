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
	"tgbot.go/config"
	"time"
)

type ConcurrentSlice struct {
	sync.RWMutex
	items []interface{}
}

type Bot struct {
	client         *tgbotapi.BotAPI
	chatToTimezone map[int64]int64
	chatToTimer    map[int64]map[string]*time.Timer
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

	log.Printf("Authorized on account %s", client.Self.UserName)

	return &Bot{client, chatToTimezone, chatToTimer}
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

	msg.Text = "TBD"
}

func (b *Bot) handleTimer(update tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	chatId := msg.ChatID
	textArray := strings.Split(update.Message.Text, " ")
	text := strings.Join(textArray[1:], " ")
	myTimer, err := NewTimer(text, chatId)
	if err != nil {
		msg.Text = err.Error()
	} else {
		b.createTimerJob(myTimer)
		msg.Text = "Timer set *successfully*!"
	}
}

func (b *Bot) createTimerJob(mTimer *Timer) {
	timeArray := strings.Split(mTimer.time, ":")
	minStr := timeArray[0]
	secStr := timeArray[1]
	min := parseInt(minStr)
	sec := parseInt(secStr)
	duration := time.Duration(min)*time.Minute + time.Duration(sec)*time.Second

	go b.createTimerChan(mTimer, duration)
}

func (b *Bot) createTimerChan(mTimer *Timer, duration time.Duration) {
	timer := time.NewTimer(duration)

	if b.chatToTimer[mTimer.chatId] == nil {
		b.chatToTimer[mTimer.chatId] = make(map[string]*time.Timer)
	}
	chatMap := b.chatToTimer[mTimer.chatId]

	shaStr := shaString(strconv.Itoa(rand.Int()))

	chatMap[shaStr] = timer

	timerEnd := <-timer.C
	delete(chatMap, shaStr)

	log.Printf("End timer for %d, %s", mTimer.chatId, timerEnd)

	msg := tgbotapi.NewMessage(mTimer.chatId, mTimer.text)
	b.client.Send(msg)

	if mTimer.repeat {
		go b.createTimerChan(mTimer, duration)
	}
}

func parseInt(str string) int64 {
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Print("error while parse str. str:" + str)
	}
	return res
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
	msg.Text += fmt.Sprintf("Active Notifications: *TBD*")
}

func (b *Bot) handleClear(update tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	res := strings.Split(update.Message.Text, " ")
	if len(res) == 2 && res[1] == "timezone" {
		b.chatToTimezone[msg.ChatID] = 0
		return
	}
	timersMap := b.chatToTimer[msg.ChatID]
	if timersMap == nil {
		return
	}
	for key, timer := range timersMap {
		timer.Stop()
		delete(timersMap, key)
	}
}