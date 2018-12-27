package cmd

import (
	"fmt"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy N",
	Short: "Copy the URL of a bookmark",
	Long:  `Copies the URL of a bookmark identified by unique number N to the clipboard.`,
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

		// Copy URL to clipboard
		clipboard.WriteAll(bm.Url.String())
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
