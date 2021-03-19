FROM golang:1.13.7-alpine3.11

# install chrome
RUN apk add --no-cache \
    libstdc++ \
    chromium \
    harfbuzz \
    nss \
    freetype \
    ttf-freefont \
    libc6-compat \
    tzdata \
    && rm -rf /var/cache/* \
    && mkdir /var/cache/apk

ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/

# set timezone UTC
RUN cp /usr/share/zoneinfo/Etc/UTC /etc/localtime \
    && echo "Etc/UTC" >  /etc/timezone

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app

COPY . .
RUN sh build.sh

CMD ["/bin/sh", "-c", "./tggo"]
