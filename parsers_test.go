package main

import (
	"github.com/PuerkitoBio/goquery"
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
