package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit N",
	Short: "Edit a bookmark",
	Long:  `Opens a text editor that allows editing a specific bookmark, identified by unique number N.`,
	Args:  IsValidBookmarkNumberArgument,
	Run: func(cmd *cobra.Command, args []string) {
		// Get bookmark number
		bookmarkNumber, _ := strconv.Atoi(args[0])

		// Load bookmarks from file
		bmks := ReadBookmarksFromFile()

		// Get existing bookmark entry
		bm, err := bmks.GetByNumber(bookmarkNumber)
		CheckError(err)

		// Edit bookmark
		EditBookmarkInEditor(&bmks, bm)

		// Save bookmarks back to file
		SaveBookmarksToFile(&bmks)

		// Print bookmark to console
		fmt.Println(FormatBookmark(bm, 0))
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
