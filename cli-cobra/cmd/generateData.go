package cmd

import (
	"cli-cobra/util"
	"fmt"

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
	MatchData := util.Scraper()

	fmt.Println(MatchData)
}
