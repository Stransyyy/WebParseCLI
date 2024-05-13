package util

import (
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type MatchDetails struct {
	Team1      string
	Team2      string
	GameStatus string
	Venue      string
	Time       string
}

func Scraper() []MatchDetails {
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

			//We parse the time
			s := time.Time{}.Format(timeStart)

			details = append(details, MatchDetails{
				Team1:      team1,
				Team2:      team2,
				GameStatus: gameStatus,
				Venue:      venue,
				Time:       s,
			})
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r.StatusCode, "\nError:", err)
	})

	c.Visit("https://www.espn.com/soccer/schedule") // Modify to the actual URL you are targeting

	return details

}
