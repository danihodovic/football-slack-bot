package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type slackMessage struct {
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	IconEmoji string `json:"icon_emoji"`
	Mrkdwn    bool   `json:"mrkdwn"`
}

func sendSlackMessage(c config, m match) {
	text := "*" + m.lastEvent().Text + " for " + m.lastEvent().Team + "*\n" + m.toString()

	msg := slackMessage{
		Channel:   c.slack.Channel,
		Username:  c.slack.Username,
		IconEmoji: c.slack.IconEmoji,
		Text:      text,
		Mrkdwn:    true,
	}

	payload, err := json.Marshal(msg)
	logErr(err)

	bufferPointer := bytes.NewBuffer(payload)
	_, err = http.Post(c.slack.URL, "application/json", bufferPointer)
	logErr(err)
}
