package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/DanNixon/voile/tui"
)

var rmCmd = &cobra.Command{
	Use:   "rm N",
	Short: "Remove a bookmark",
	Long:  `Removes a bookmark identified by unique number N from the library.`,
	Args:  IsValidBookmarkNumberArgument,
	Run: func(cmd *cobra.Command, args []string) {
		// Get bookmark number
		bookmarkNumber, _ := strconv.Atoi(args[0])

		// Load bookmarks from file
		bmks := ReadBookmarksFromFile()

		// Get bookmark and print to console
		bm, err := bmks.GetByNumber(bookmarkNumber)
		CheckError(err)
		fmt.Println(FormatBookmark(bm, 0))

		// Determine if bookmark should be deleted
		rm, _ := cmd.Flags().GetBool(ForceFlagName)
		if !rm {
			rm, _ = tui.Confirm("Really remove bookmark?")
		}

		if rm {
			// Remove bookmark
			err := bmks.DeleteByNumber(bookmarkNumber)
			CheckError(err)
		} else {
			fmt.Println("Bookmark not removed.")
		}

		// Save bookmarks back to file
		SaveBookmarksToFile(&bmks)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	rmCmd.Flags().Bool(ForceFlagName, false, "Remove without confirmation")
}
