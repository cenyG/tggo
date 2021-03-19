## TgBot
This is educational project to learn *GoLang* and also *Telegram API*.   

Bot have following commands:
```
/help      - Help
/status    - Show your active timers/notificators
/screen    - Make instant screen
/every     - Make screen every `m(min)`, `h(hour)`, `d(day)`. Example: `/every 5m google.com
```

##### Run dev:
```
./dev.sh
```

##### Run Docker:
```
docker build .
docker run -d {IMAGE_NAME}
```

You will need to create `.env` file contains `BOT_TOKEN` and `PROXY`(optional) to fill.
