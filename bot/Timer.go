package bot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Timer struct {
	duration time.Duration
	text     string
	repeat   bool
	chatId   int64
}

func (t Timer) String() string {
	return fmt.Sprintf(`Timer: *** duration: %s, text: %s, repeat: %t, chatId: %d ***`, t.duration.String(), t.text, t.repeat, t.chatId)
}

func NewTimer(textArray []string, chatId int64) (*Timer, error) {
	if len(textArray) == 0 {
		return nil, errors.New(`bad format`)
	}
	if textArray[0] == "every" { // every mm:ss text
		mTime := textArray[1]
		text := strings.Join(textArray[2:], " ")
		if text == "" {
			text = defaultText
		}
		duration, err := parseMinSec(mTime)
		if err != nil {
			return nil, err
		}
		return &Timer{duration, text, true, chatId}, nil
	} else { //mm:ss text
		mTime := textArray[0]
		text := strings.Join(textArray[1:], " ")
		if text == "" {
			text = defaultText
		}
		duration, err := parseMinSec(mTime)
		if err != nil {
			return nil, err
		}
		return &Timer{duration, text, false, chatId}, nil
	}
}

func parseMinSec(text string) (time.Duration, error) {
	res := strings.Split(text, ":")

	if len(res) == 2 {
		if len(res[0]) > 1 && len(res[1]) == 2 {
			min, err1 := strconv.ParseInt(res[0], 10, 64)
			sec, err2 := strconv.ParseInt(res[1], 10, 64)

			if min < 0 || min > maxMin || sec < 0 || sec > maxSec || err1 != nil || err2 != nil {
				return 0, errors.New("bad min and sec format")
			}
			duration := time.Duration(min)*time.Minute + time.Duration(sec)*time.Second

			return duration, nil
		}
	}

	return 0, errors.New("bad min and sec format")
}
