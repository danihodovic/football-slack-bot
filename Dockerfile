FROM golang:1.6

RUN go get \
	github.com/PuerkitoBio/goquery \
	gopkg.in/redis.v4

COPY . /app

WORKDIR /app

CMD go build -o /usr/local/bin/app && app --config config.json

