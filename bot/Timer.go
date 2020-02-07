package bot

import (
	"errors"
	"strings"
)

type Timer struct {
	time   string
	text   string
	repeat bool
	chatId int64
}

const defaultText = "⏰⏰⏰ Alarmed ! ! ! ⏰⏰⏰"

func NewTimer(text string, chatId int64) (*Timer, error) {
	res := strings.Split(text, " ")

	if res[0] == "every" { // every mm:ss text
		mTime := res[1]
		text := strings.Join(res[2:], " ")
		if text == "" {
			text = defaultText
		}
		err := validateMinSec(mTime)
		if err != nil {
			return nil, err
		}
		return &Timer{mTime, text, true, chatId}, nil
	} else { //mm:ss text
		mTime := res[0]
		text := strings.Join(res[1:], " ")
		if text == "" {
			text = defaultText
		}
		err := validateMinSec(mTime)
		if err != nil {
			return nil, err
		}
		return &Timer{mTime, text, false, chatId}, nil
	}
}

func validateMinSec(text string) error {
	res := strings.Split(text, ":")
	if len(res) == 2 {
		if len(res[0]) > 1 && len(res[1]) == 2 {
			return nil
		}
	}

	return errors.New("bad min and sec format ")
}
