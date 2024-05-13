package util

import (
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var (
	layout = "3:04 PM"
)

// MatchDetails holds the details of a match
type MatchDetails struct {
	Team1      string
	Team2      string
	GameStatus string
	Venue      string
	Time       string
}

// TimeResult holds the parsed time and any error message
type TimeResult struct {
	ParsedTime string
	ErrorMsg   string
}

// ParseAndValidateTime parses and validates the given time string and adjusts it to the specified time zone
func ParseAndValidateTime(value string, timeZone string) TimeResult {
	value = strings.TrimSpace(value)

	t, err := time.Parse(layout, value)
	if err != nil {
		return TimeResult{
			ParsedTime: "",
			ErrorMsg:   "Time cannot be parsed",
		}
	}

	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return TimeResult{
			ParsedTime: "",
			ErrorMsg:   err.Error(),
		}
	}

	t = t.In(loc)
	return TimeResult{
		ParsedTime: t.Format(layout),
		ErrorMsg:   "",
	}
}

// Scraper scrapes match details from the specified website
func Scraper(timeZone string) []MatchDetails {
	c := colly.NewCollector(
		colly.AllowedDomains("espn.com", "www.espn.com"), // Assuming ESPN domain, adjust as necessary
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"),
	)

	var details []MatchDetails

	details = append(details, MatchDetails{
		Team1:      "Team 1",
		Team2:      "Team 2",
		GameStatus: "Game Status",
		Venue:      "Venue",
		Time:       "Time",
	})

	c.OnHTML("tbody.Table__TBODY", func(e *colly.HTMLElement) {
		e.ForEach("tr.Table__TR--sm", func(_ int, el *colly.HTMLElement) {
			team1 := el.ChildText("span.Table__Team.away a.AnchorLink")
			team2 := el.ChildText("span.Table__Team a.AnchorLink:last-child")
			gameStatus := el.ChildText("span.gameNote")
			timeStart := el.ChildText("td.date__col a.AnchorLink")
			venue := el.ChildText("td.venue__col div")
			gameStatusV2 := el.ChildText("td.teams__col a.AnchorLink")

			if strings.TrimSpace(gameStatusV2) == "FT" {
				gameStatus = "Full Time"
			}

			// We need to convert the time to a more readable format
			timeResult := ParseAndValidateTime(timeStart, timeZone)
			if timeResult.ErrorMsg != "" {
				log.Println("Error parsing time:", timeResult.ErrorMsg)
				timeResult.ParsedTime = "FT"
			}

			// Append the data to the slice
			details = append(details, MatchDetails{
				Team1:      team1,
				Team2:      team2,
				GameStatus: gameStatus,
				Venue:      venue,
				Time:       timeResult.ParsedTime,
			})
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r.StatusCode, "\nError:", err)
	})

	c.Visit("https://www.espn.com/soccer/schedule") // Modify to the actual URL you are targeting

	return details
}
