package main

import (
	"flag"
	"gopkg.in/redis.v4"
	"log"
	"time"
)

func logErr(err error) {
	if err != nil {
		panic(err)
	}
}

func run(client *redis.Client, cfg config) {
	matches := parseESPN("http://www.espnfc.com")

	for _, m := range matches {
		if !relevantTeam(cfg, m.HomeTeam) || !relevantTeam(cfg, m.AwayTeam) {
			continue
		}

		oldM, err := getMatch(client, m.toKey())
		if err != nil && err != redis.Nil {
			panic(err)
		}

		if err == redis.Nil {
			setMatch(client, m)
			log.Println("Redis key zero, setting", m)
			continue
		}

		// We have new events
		if relevantEvent(cfg, m) && len(m.MatchEvents) > len(oldM.MatchEvents) {
			log.Println("New event!", m.lastEvent().Text, m.toString())
			log.Println(m.MatchEvents, oldM.MatchEvents)
			sendSlackMessage(cfg, m)
			setMatch(client, m)
		}
	}
}

func main() {
	filterFile := flag.String("config", "", "The JSON configuration file")
	flag.Parse()
	cfg := parseConfig(*filterFile)

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	for {
		log.Println("Parsing...")
		run(client, cfg)
		time.Sleep(cfg.interval)
	}
}
