package bot

import "testing"

func TestScreener_MakeScreen(t *testing.T) {
	screener := &Screener{url: `https://pass.rzd.ru/tickets/public/ru?STRUCTURE_ID=704&refererPageId=4819&layer_name=e3-route&tfl=3&st0=%D0%9C%D0%BE%D1%81%D0%BA%D0%B2%D0%B0&code0=2000000&dt0=16.02.2020&st1=%D0%A1%D0%B0%D0%BD%D0%BA%D1%82-%D0%9F%D0%B5%D1%82%D0%B5%D1%80%D0%B1%D1%83%D1%80%D0%B3&code1=2004000&checkSeats=1`}

	screener.MakeScreen()
}
