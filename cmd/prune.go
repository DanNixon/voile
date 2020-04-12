package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/DanNixon/voile/tui"
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Bring bookmarks up to date",
	Long:  `Picks an old bookmark to determine if it still relevant.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load bookmarks from file
		bmks := ReadBookmarksFromFile()
		sort.SliceStable(bmks.Bookmarks, func(i, j int) bool {
			return bmks.Bookmarks[i].LastUpdated.Before(bmks.Bookmarks[j].LastUpdated)
		})

		// Get existing bookmark entry
		bmNumber := bmks.Bookmarks[0].Number

		// Resort bookmarks in natural order
		sort.Sort(&bmks)

		bm, err := bmks.GetByNumber(bmNumber)
		CheckError(err)

		// Print bookmark to console
		fmt.Println(FormatBookmark(bm, 0))

		// Determine what to do with bookmark
		result, err := tui.Option("Action", []tui.MultiChoiceOption{{"k", "keep"}, {"e", "edit"}, {"d", "delete"}})
		CheckError(err)
		if result == "e" {
			EditBookmarkInEditor(&bmks, bm)
			fmt.Println(FormatBookmark(bm, 0))
		} else if result == "k" {
			bm.MarkUpdated()
		} else if result == "d" {
			bmks.DeleteByNumber(bmNumber)
		}

		// Save bookmarks back to file
		SaveBookmarksToFile(&bmks)
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
}
