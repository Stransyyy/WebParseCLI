package cmd

import (
	"cli-cobra/util"
	"log"

	"github.com/spf13/cobra"
)

var loadMatches = &cobra.Command{
	Use:   "fulbo",
	Short: "Load all the matches from today",
	Long:  `This command will load all the matches from today and store them in a CSV file, or the type of file you specify.`,
	Run:   loadGames,
}

func init() {
	rootCmd.AddCommand(loadMatches)

}

func loadGames(cmd *cobra.Command, args []string) {
	MatchData := util.Scraper("America/Chicago")

	for _, match := range MatchData {
		log.Printf("Match: %s vs %s, Status: %s, Venue: %s, Time: %s, Score: %s\n\n", match.Team1, match.Team2, match.GameStatus, match.Venue, match.Time, match.Score)
	}

	//fmt.Println(MatchData)
}
