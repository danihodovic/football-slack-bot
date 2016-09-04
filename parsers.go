package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type matchEvent struct {
	Minute    string `json:"minute"`
	EventType string `json:"event_type"`
	Team      string `json:"team"`
	Text      string `json:"text"`
}

func (m matchEvent) sortableMinute() int {
	var minute int
	if m.Minute == "HT" {
		return 45
		//Sort extra time to the previous 45
	}
	if strings.Contains(m.Minute, "+") {
		return 45
	}
	var err error
	minute, err = strconv.Atoi(m.Minute)
	logErr(err)
	return minute
}

type byMinute []matchEvent

func (s byMinute) Len() int {
	return len(s)
}
func (s byMinute) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byMinute) Less(i, j int) bool {
	return s[i].sortableMinute() < s[j].sortableMinute()
}

type match struct {
	TimeCurrent   string       `json:"time_current"`
	HomeTeam      string       `json:"home_team"`
	AwayTeam      string       `json:"away_team"`
	HomeTeamGoals int          `json:"home_team_goals"`
	AwayTeamGoals int          `json:"away_team_goals"`
	MatchID       string       `json:"match_id"`
	MatchEvents   []matchEvent `json:"match_events"`
}

func (m *match) toKey() string {
	return fmt.Sprintf("%s:%s", m.HomeTeam, m.AwayTeam)
}

func (m *match) toJSON() []byte {
	bytes, err := json.Marshal(m)
	logErr(err)
	return bytes
}

func (m *match) lastEvent() *matchEvent {
	if len(m.MatchEvents) > 0 {
		return &m.MatchEvents[len(m.MatchEvents)-1]
	}
	return nil
}

func (m *match) toString() string {
	return fmt.Sprintf("%s' %s %d:%d %s",
		m.TimeCurrent, m.HomeTeam, m.HomeTeamGoals, m.AwayTeamGoals, m.AwayTeam)
}

func (m *match) ttl() time.Duration {
	return time.Duration(120 * time.Minute)
	// var currentMinute int
	// if m.TimeCurrent == "HT" {
	// currentMinute = 45
	// } else {
	// var err error
	// currentMinute, err = strconv.Atoi(m.TimeCurrent)
	// if err != nil {
	// log.Println("Error parsing time remaining for: ", m, "Setting to 120")
	// return time.Duration(120)
	// }
	// }

	// maxTimeRemaining := 120 - currentMinute
	// return time.Duration(maxTimeRemaining) * time.Second
}

func parseESPN() []match {
	url := "http://www.espnfc.com/scores"
	doc, err := goquery.NewDocument(url)
	logErr(err)

	var matches []match

	doc.Find("div.scorebox-container.live").Each(func(idx int, scoreDiv *goquery.Selection) {
		href, exists := scoreDiv.Children().Eq(1).Attr("href")
		if !exists {
			log.Fatalln("No href found")
		}
		matchIDRegex := regexp.MustCompile(`gameId=(\d+)`)
		matchID := matchIDRegex.FindStringSubmatch(href)[1]

		matchURL := "http://www.espnfc.com/match?gameId=%s"
		matchURL = fmt.Sprintf(matchURL, matchID)

		matchDoc, err := goquery.NewDocument(matchURL)
		logErr(err)

		homeTeam := matchDoc.Find("div.away span.short-name").Text()
		awayTeam := matchDoc.Find("div.home span.short-name").Text()

		homeTeamScoreStr := strings.TrimSpace(matchDoc.Find("div.away span.score").Text())
		homeTeamScore, err := strconv.Atoi(homeTeamScoreStr)
		// Sometimes the str is "" which means that the result is 0.
		if err != nil {
			homeTeamScore = 0
		}

		awayTeamScoreStr := strings.TrimSpace(matchDoc.Find("div.home span.score").Text())
		awayTeamScore, err := strconv.Atoi(awayTeamScoreStr)
		if err != nil {
			awayTeamScore = 0
		}

		currentTime := matchDoc.Find("span.game-time").Text()
		if strings.Contains(currentTime, "'") {
			currentTime = currentTime[0:2]
		}

		m := match{
			TimeCurrent:   currentTime,
			HomeTeam:      homeTeam,
			AwayTeam:      awayTeam,
			HomeTeamGoals: homeTeamScore,
			AwayTeamGoals: awayTeamScore,
			MatchID:       matchID,
		}

		events := parseESPNMatchDetails(matchDoc, m)
		m.MatchEvents = events

		// The time parsed at the top of the page and the time in an even <span> are different, with
		// the event span being newer -.-
		if m.lastEvent() != nil {
			m.TimeCurrent = m.lastEvent().Minute
		}

		matches = append(matches, m)
	})

	return matches
}

func parseESPNMatchDetails(doc *goquery.Document, m match) []matchEvent {
	var events []matchEvent

	doc.Find("li[data-time]").Each(func(idx int, li *goquery.Selection) {
		minuteStr, _ := li.Attr("data-time")
		if minuteStr == "KO" {
			return
		}

		homeEventsUl := li.Find("ul[data-event-home-away='home']")
		if homeEventsUl.Length() > 0 {
			ev := parseTeamEvent(homeEventsUl, m.HomeTeam, minuteStr)
			events = append(events, ev)
		}

		awayEventsUl := li.Find("ul[data-event-home-away='away']")
		if awayEventsUl.Length() > 0 {
			ev := parseTeamEvent(awayEventsUl, m.AwayTeam, minuteStr)
			events = append(events, ev)
		}
	})

	sort.Sort(byMinute(events))
	return events
}

func parseTeamEvent(homeEvent *goquery.Selection, team, minute string) matchEvent {
	eventType, _ := homeEvent.Find("li").Attr("data-events-type")
	eventType = strings.Replace(eventType, "-", " ", -1)
	eventType = strings.ToLower(eventType)

	text := strings.TrimSpace(homeEvent.Find("div.detail").Text())

	ev := matchEvent{
		Minute:    minute,
		EventType: eventType,
		Team:      team,
		Text:      text,
	}

	return ev
}
