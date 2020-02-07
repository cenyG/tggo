package bot

import (
	"errors"
	"strings"
)

type Notification struct {
	time   string
	text   string
	repeat bool
	chatId int64
}

func NewNotification(text string, chatId int64) (*Notification, error) {
	res := strings.Split(text, " ")

	if len(res) == 2 { // HH:mm text
		mTime := res[0]
		text := res[1]

		err := validateHHmm(mTime)
		if err != nil {
			return nil, err
		}

		return &Notification{
			mTime,
			text,
			false,
			chatId,
		}, nil
	} else if len(res) == 3 {
		if res[0] == "every" { //every HH:mm text
			mTime := res[1]
			text := res[2]

			err := validateHHmm(mTime)
			if err != nil {
				return nil, err
			}

			return &Notification{
				mTime,
				text,
				true,
				chatId,
			}, nil

		} else { // DD/MM/YYYY HH:mm text
			mDay := res[0]
			mTime := res[1]
			text := res[2]

			err1, err2 := validateDDMMYYYY(mDay), validateHHmm(mTime)
			if err1 != nil || err2 != nil {
				return nil, errors.New("bad date/time format")
			}

			mDateTime := strings.Join([]string{mDay, mTime}, " ")
			return &Notification{
				mDateTime,
				text,
				false,
				chatId,
			}, nil

		}
	}
	return nil, errors.New("bad text format")
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
