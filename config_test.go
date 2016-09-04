package main

import (
	"reflect"
	"testing"
	"time"
)

func TestParseConfig(t *testing.T) {
	cfg := parseConfig("config-sample.json")

	expectedSlackConfig := slackConfig{
		URL:       "https://hooks.slack.com/services/foo/bar/baz",
		Channel:   "#test",
		Username:  "test",
		IconEmoji: "test",
	}

	expectedConfig := config{
		interval: time.Duration(30) * time.Second,
		teams: map[string]bool{
			"team to notify on 1": true,
			"team to notify on 2": true,
		},
		events: map[string]bool{
			"yellow card":  true,
			"red card":     true,
			"goal":         true,
			"substitution": true,
		},
		slack: expectedSlackConfig,
	}

	if !reflect.DeepEqual(expectedConfig, cfg) {
		t.Fatalf("Config struct not eql to what was expected")
	}
}
