package cmd

import (
	"fmt"
	"strconv"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open N",
	Short: "Open a bookmark",
	Long:  `Opens a bookmark identified by unique number N in a web browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get bookmark number
		bookmarkNumber, _ := strconv.Atoi(args[0])

		// Load bookmarks from file
		bmks := ReadBookmarksFromFile()

		// Get bookmark entry
		bm, err := bmks.GetByNumber(bookmarkNumber)
		CheckError(err)

		// Print bookmark to console
		fmt.Println(FormatBookmark(bm, 0))

		// Open URL in browser
		open.Run(bm.Url.String())
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
