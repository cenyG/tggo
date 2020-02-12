##TgBot
This is educational project to learn *GoLang* and also *Telegram API*.   

Bot have following commands:
```
/help      - Help
/status    - Show your active timers/notificators
/timezone  - Set your timezone
/set       - Create notification: HH:mm, DD/MM/YYYY
/timer     - Create timer: mm:SS
/clear     - Clear all timers/notifications
```

#####Run dev:
```
./dev.sh
```

#####Run Docker:
```
docker build .
docker run -d {IMAGE_NAME}
```

`.env` file contains `BOT_TOKEN` and `PROXY`(optional) to fill.