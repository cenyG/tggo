package bot

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func timezoneParse(timezone string) (int64, error) {
	if len(timezone) > 3 {
		return 0, errors.New("bad timezone")
	} else {
		sign := timezone[0]
		if sign != '+' && sign != '-' {
			return 0, errors.New("bad timezone")
		}
		number, err := strconv.ParseInt(timezone[1:], 10, 64)
		if err != nil {
			return 0, err
		}
		if sign == '-' {
			number = -1 * number
		}

		return number, nil
	}
}

func getHttpClient() *http.Client {
	proxy := os.Getenv("PROXY")
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
