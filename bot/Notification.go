package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type Notification struct {
	duration time.Duration
	text     string
	repeat   bool
	chatId   int64
}

func (n Notification) String() string {
	return fmt.Sprintf(`Notification: *** duration: %d, text: %s, repeat: %t, chatId: %d ***`, n.duration, n.text, n.repeat, n.chatId)
}

func NewNotification(textArray []string, chatId int64, utcDiff int64) (*Notification, error) {
	first := textArray[0]
	second := ""
	if len(textArray) > 1 {
		second = textArray[1]
	}

	var timeString string
	var text string
	var repeat bool

	if first == "every" {
		err := validateHHmm(second)
		if err != nil {
			return nil, err
		}
		timeString = second
		text = strings.Join(textArray[2:], " ")
		repeat = true

	} else if validateDDMMYYYY(first) == nil { // HH:mm text
		err := validateHHmm(second)
		if err != nil {
			return nil, err
		}
		timeString = strings.Join([]string{first, second}, " ")
		text = strings.Join(textArray[2:], " ")

	} else if validateHHmm(first) == nil {
		timeString = first
		text = strings.Join(textArray[1:], " ")

	} else {
		return nil, errors.New("bad text format")
	}

	if text == "" {
		text = defaultText
	}

	duration, err := parseTimeString(timeString, utcDiff)
	if err != nil {
		return nil, err
	}

	return &Notification{
		duration,
		text,
		repeat,
		chatId,
	}, nil
}

func parseTimeString(timeString string, utcDiff int64) (time.Duration, error) {
	timeNow := time.Now().UTC()
	var duration time.Duration

	if len(timeString) == 16 {
		mTime, err := time.Parse("02/01/2006 15:04", timeString)
		if err != nil {
			return 0, errors.New("wrong date format")
		}
		mTime = mTime.Add(time.Duration(-1*utcDiff) * time.Hour)

		duration = mTime.Sub(timeNow)
		if duration < 0 {
			log.Printf("Duration %s, mTime: %s", duration, mTime)
			return 0, errors.New("wrong date format")
		}

	} else if len(timeString) == 5 {
		mTime, err := time.Parse("15:04", timeString)
		if err != nil {
			return 0, errors.New("wrong date format")
		}
		mTime = mTime.Add(time.Duration(-1*utcDiff) * time.Hour)

		timeTo := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), mTime.Hour(), mTime.Minute(), 0, 0, time.UTC)

		if timeTo.Unix() < timeNow.Unix() {
			timeTo.AddDate(0, 0, 1)
		}

		duration = timeTo.Sub(timeNow)

	} else {
		log.Printf("Wrong date: %s", timeString)
		return 0, errors.New("wrong date format")
	}

	return duration, nil
}

func validateHHmm(text string) error {
	res := strings.Split(text, ":")
	if len(res) == 2 {
		if len(res[0]) == 2 && len(res[1]) == 2 {
			return nil
		}
	}
	return errors.New("bad hour and minute format")
}

func validateDDMMYYYY(text string) error {
	res := strings.Split(text, "/")
	if len(res) == 3 {
		if len(res[0]) == 2 && len(res[1]) == 2 && len(res[2]) == 4 {
			return nil
		}
	}
	return errors.New("bad day, month and year format")
}
