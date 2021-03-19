package bot

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"strconv"
	"time"
)

type Screener struct {
	duration time.Duration
	url      string
	chatId   int64
}

func (s *Screener) String() string {
	return fmt.Sprintf(`Screener: *** duration: %s, url: %s, chatId: %d ***`, s.duration.String(), s.url, s.chatId)
}

func NewScreener(chatId int64, time string, url string) (*Screener, error) {
	duration, err := parseTimeString(time)
	if err != nil {
		return nil, err
	}

	return &Screener{
		*duration,
		parseUrl(url),
		chatId,
	}, nil
}

func (s *Screener) MakeScreen() (string, error) {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), chromedp.NoSandbox, chromedp.Headless)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx,
		fullScreenshot(s.url, 90, &buf),
	); err != nil {
		return ``, err
	}

	screenFileNamePattern := time.Now().UTC().Format(`2006-01-02T15:04_`) + strconv.FormatInt(s.chatId, 10) + "_*" + ".png"
	tmpFile, err := ioutil.TempFile(``, screenFileNamePattern)
	if err != nil {
		return ``, err
	}
	defer closeFile(tmpFile)

	if _, err := tmpFile.Write(buf); err != nil {
		return ``, err
	}

	return tmpFile.Name(), nil
}
