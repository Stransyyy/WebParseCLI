package util

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var (
	layout  = "3:04 PM" // Parsing 12-hour clock format with AM/PM
	pattern = `^\d+\s*-\s*\d+$`
)

// MatchDetails holds the details of a match
type MatchDetails struct {
	Team1      string
	Team2      string
	GameStatus string
	Venue      string
	Time       string
	Score      string
}

// TimeResult holds the parsed time and any error message
type TimeResult struct {
	ParsedTime string
	ErrorMsg   string
}

// ParseAndValidateTime parses and validates the given time string and adjusts it to the specified time zone
func ParseAndValidateTime(value string, timeZone string) TimeResult {
	// Clean the input time string
	value = strings.TrimSpace(value)

	if value == "" {
		return TimeResult{
			ParsedTime: "",
			ErrorMsg:   "Match time is empty or the match has already been played",
		}
	}

	// Debug print the raw time string
	fmt.Println("Raw time string:", value)

	// Parse the time string according to the layout
	t, err := time.Parse(layout, value)
	if err != nil {
		fmt.Println("Parsing error:", err)
		return TimeResult{
			ParsedTime: "",
			ErrorMsg:   "Time cannot be parsed",
		}
	}

	// Load the specified time zone
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return TimeResult{
			ParsedTime: "",
			ErrorMsg:   err.Error(),
		}
	}

	// Adjust the parsed time to the specified time zone
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
	}) // make this the header of the csv file

	c.OnHTML("tbody.Table__TBODY", func(e *colly.HTMLElement) {
		e.ForEach("tr.Table__TR--sm", func(_ int, el *colly.HTMLElement) {
			team1 := el.ChildText("span.Table__Team.away a.AnchorLink")
			team2 := el.ChildText("span.Table__Team a.AnchorLink")                          //
			gameStatus := el.ChildText("span.gameNote")                                     // Game status, like ie: FT, HT, 1st Leg, etc
			timeStatus := el.ChildText("td.date__col a.AnchorLink")                         // Time of start of the game
			timeStatus2 := el.ChildText("td.date__col a.Schedule__liveLink clr-brand-ESPN") // check for LIVE game Time
			venue := el.ChildText("td.venue__col div")
			gameStatusV2 := el.ChildText("td.teams__col a.AnchorLink")
			score := el.ChildText("td.colspan__col.Table__TD a.AnchorLink.at")

			reg := regexp.MustCompile(pattern)

			if timeStatus == timeStatus2 {
				timeStatus = timeStatus2
			}

			if gameStatus == "" {
				gameStatus = gameStatusV2
			}

			// Append the data to the slice
			if reg.MatchString(score) {

				details = append(details, MatchDetails{
					Team1:      team1,
					Team2:      team2,
					GameStatus: gameStatus,
					Venue:      venue,
					Time:       timeStatus,
					Score:      score,
				})

			} else {
				details = append(details, MatchDetails{
					Team1:      team1,
					Team2:      team2,
					GameStatus: gameStatus,
					Venue:      venue,
					Time:       timeStatus,
					Score:      "Not Available",
				})
			}
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r.StatusCode, "\nError:", err)
	})

	c.Visit("https://www.espn.com/soccer/schedule") // Modify to the actual URL you are targeting

	return details
}
