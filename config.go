package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"
)

type slackConfig struct {
	URL       string
	Channel   string
	Username  string
	IconEmoji string
}

type config struct {
	interval time.Duration
	teams    map[string]bool
	events   map[string]bool
	slack    slackConfig
}

func parseConfig(filename string) config {
	cfg := config{
		teams:  make(map[string]bool),
		events: make(map[string]bool),
	}

	b, err := ioutil.ReadFile(filename)
	logErr(err)

	var configJSON map[string]interface{}
	err = json.Unmarshal(b, &configJSON)
	logErr(err)

	for _, v := range configJSON["teams"].([]interface{}) {
		str := v.(string)
		str = strings.ToLower(str)
		cfg.teams[str] = true
	}

	for _, v := range configJSON["events"].([]interface{}) {
		str := v.(string)
		str = strings.ToLower(str)
		cfg.events[str] = true
	}

	slackConfig := configJSON["slack"].(map[string]interface{})
	cfg.slack.URL = slackConfig["url"].(string)
	cfg.slack.Channel = slackConfig["channel"].(string)
	cfg.slack.Username = slackConfig["username"].(string)
	cfg.slack.IconEmoji = slackConfig["iconEmoji"].(string)

	cfg.interval = time.Duration(configJSON["interval"].(float64)) * time.Second

	return cfg
}

func relevantTeam(cfg config, team string) bool {
	_, relevant := cfg.teams[strings.ToLower(team)]
	return relevant
}

func relevantEventType(cfg config, eventType string) bool {
	_, relevant := cfg.events[strings.ToLower(eventType)]
	return relevant
}

func relevantEvent(cfg config, m match) bool {
	if m.lastEvent() == nil {
		return false
	}
	team := (relevantTeam(cfg, m.HomeTeam) || relevantTeam(cfg, m.AwayTeam))
	eventType := relevantEventType(cfg, m.lastEvent().EventType)
	return team && eventType
}
