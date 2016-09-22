package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestParseESPNMatchDetails(t *testing.T) {
	fileReader, err := os.Open("test-samples/example-match-details-gameId-453699.html")
	if err != nil {
		t.Fail()
	}

	doc, err := goquery.NewDocumentFromReader(fileReader)
	if err != nil {
		t.Fail()
	}

	m := match{
		HomeTeam: "home",
		AwayTeam: "away",
	}

	events := parseESPNMatchDetails(doc, m)

	expectedEvents := []matchEvent{
		matchEvent{Minute: "8", EventType: "goal", Team: "home", Text: "Harry Beautyman Goal"},
		matchEvent{Minute: "13", EventType: "goal", Team: "home", Text: "Alex Revell Goal"},
		matchEvent{Minute: "28", EventType: "goal", Team: "home", Text: "Matthew Taylor Goal"},
		matchEvent{Minute: "38", EventType: "goal", Team: "away", Text: "Dean Bowditch Goal"},
		matchEvent{Minute: "43", EventType: "yellow card", Team: "away", Text: "Samir Carruthers Yellow Card"},
		matchEvent{Minute: "44", EventType: "yellow card", Team: "away", Text: "George Baldock Yellow Card"},
		matchEvent{Minute: "59", EventType: "substitution", Team: "away", Text: "Kieran Agard Substitution"},
		matchEvent{Minute: "59", EventType: "substitution", Team: "away", Text: "Ben Reeves Substitution"},
		matchEvent{Minute: "68", EventType: "substitution", Team: "home", Text: "Alfie Potter Substitution"},
		matchEvent{Minute: "68", EventType: "substitution", Team: "home", Text: "Sam Hoskins Substitution"},
		matchEvent{Minute: "71", EventType: "substitution", Team: "away", Text: "Daniel Powell Substitution"},
	}

	if !reflect.DeepEqual(expectedEvents, events) {
		t.Log("expected events and events not equal")
	}
}

func TestParseESPN(t *testing.T) {
	http.HandleFunc("/scores", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./test-samples/espn_example_1_scores.html")
	})

	http.HandleFunc("/match", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query()["gameId"][0] == "466482" {
			http.ServeFile(w, r, "./test-samples/espn_example_1_live_466482.html")
		}
	})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	matches := parseESPN("http://localhost:8080")
	var match466482 match

	for _, match := range matches {
		if match.MatchID == "466482" {
			match466482 = match
		}
	}

	expectedMatch := match{
		TimeCurrent:   "54",
		HomeTeam:      "Twente Enschede",
		AwayTeam:      "FC Utrecht",
		HomeTeamGoals: 1,
		AwayTeamGoals: 1,
		MatchID:       "466482",
		MatchEvents: []matchEvent{
			matchEvent{Minute: "17", EventType: "goal", Team: "FC Utrecht", Text: "Richairo Zivkovic Goal"},
			matchEvent{Minute: "54", EventType: "goal", Team: "Twente Enschede", Text: "Enes Unal Goal"},
		},
	}

	if !reflect.DeepEqual(match466482, expectedMatch) {
		t.Fail()
	}
}
