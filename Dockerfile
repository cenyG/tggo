FROM golang:1.13.7-alpine3.11

WORKDIR /app
RUN apk add --no-cache \
        libc6-compat
RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Etc/UTC /etc/localtime
RUN echo "Etc/UTC" >  /etc/timezone

COPY . .
RUN sh build.sh

ENTRYPOINT ["/bin/sh", "-c", "./tggo"]
