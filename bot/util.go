package bot

import (
	"context"
	"errors"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"tgbot/config"
	"time"
)

func fullScreenshot(url string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(5 * time.Second)
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

// input string should be like 5m, 2h, 1d
// returns time converted to time.Duration
func parseTimeString(str string) (*time.Duration, error) {
	runeStr := []rune(str)
	if len(runeStr) < 2 {
		return nil, errors.New("wrong time format")
	}

	intervalStr := runeStr[0 : len(runeStr)-1]
	intervalSign := runeStr[len(runeStr)-1]

	interval, err := strconv.ParseInt(string(intervalStr), 10, 64)
	if err != nil {
		return nil, err
	}

	var res time.Duration

	if intervalSign == 'm' {
		res = time.Duration(interval) * time.Minute
	} else if intervalSign == 'h' {
		res = time.Duration(interval) * time.Hour
	} else if intervalSign == 'd' {
		res = time.Duration(interval) * time.Hour * 24
	} else {
		return nil, errors.New("wrong time format")
	}

	return &res, nil
}

func getHttpClient() *http.Client {
	proxy := config.GetProxy()
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			log.Println(err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

		httpClient := &http.Client{
			Transport: transport,
		}
		return httpClient
	}
	return &http.Client{}
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Println(`error can't close file'`, err)
	}
}

func removeFile(filePath string) {
	if err := os.Remove(filePath); err != nil {
		log.Println(`error while file delete`, err)
	}
}

func parseUrl(url string) string {
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		url = "https://" + url
	}

	return url
}