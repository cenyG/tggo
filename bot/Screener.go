package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Screener struct {
	duration time.Duration
	url      string
	repeat   bool
	chatId   int64
}

func (s Screener) String() string {
	return fmt.Sprintf(`Screener: *** duration: %s, url: %s, repeat: %t, chatId: %d ***`, s.duration.String(), s.url, s.repeat, s.chatId)
}

func NewScreener(textArray []string, chatId int64, utcDiff int64) (*Screener, error) {
	first := textArray[0]
	second := ""
	if len(textArray) > 1 {
		second = textArray[1]
	}

	var timeString string
	var url string
	var repeat bool

	if first == "every" {
		err := validateHHmm(second)
		if err != nil {
			return nil, err
		}
		timeString = second
		if len(textArray) > 3 {
			return nil, errors.New(`bad format`)
		}
		url = textArray[2]
		repeat = true

	} else if validateDDMMYYYY(first) == nil { // DD/MM/YYYY HH:mm url
		err := validateHHmm(second)
		if err != nil {
			return nil, err
		}
		timeString = strings.Join([]string{first, second}, " ")
		if len(textArray) > 3 {
			return nil, errors.New(`bad format`)
		}
		url = textArray[2]

	} else if validateHHmm(first) == nil { //HH:mm url
		timeString = first
		if len(textArray) > 2 {
			return nil, errors.New(`bad format`)
		}
		url = textArray[1]

	} else {
		return nil, errors.New("bad text format")
	}

	if url == "" {
		return nil, errors.New(`no url`)
	}

	duration, err := parseTimeString(timeString, utcDiff)
	if err != nil {
		return nil, err
	}

	return &Screener{
		duration,
		url,
		repeat,
		chatId,
	}, nil
}

func (s *Screener) MakeScreen() (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
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

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Println(`error can't close file'`, err)
	}
}

func fullScreenshot(url string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(20 * time.Second)
			return nil
		}),
		//chromedp.WaitReady(`html`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
