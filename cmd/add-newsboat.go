package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addNewsboatCmd = &cobra.Command{
	Use:   "add-newsboat URL TITLE DESCRIPTION (FEED TITLE)",
	Short: "Add a new bookmark from Newsboat",
	Long:  `Adds a new bookmark using commands passed by Newboats bookmarking system.`,
	Args:  cobra.RangeArgs(3, 4),
	Run: func(cmd *cobra.Command, args []string) {
		// Load bookmarks from file
		bmks := ReadBookmarksFromFile()

		// Create new bookmark entry
		bm := bmks.NewEntry()

		// Set URL
		err := bm.Url.Parse(args[0])
		CheckError(err)

		bm.Name = args[1]
		bm.Description = args[2]

		// Edit in editor if requested
		editFlag, _ := cmd.Flags().GetBool(EditFlagName)
		if editFlag {
			// Validate the bookmarks before opening editor
			err = bmks.Verify()
			CheckError(err)

			EditBookmarkInEditor(&bmks, bm)
		}

		// Save bookmarks back to file
		SaveBookmarksToFile(&bmks)

		// Print bookmark to console
		fmt.Println(FormatBookmark(bm, 0))
	},
}

func init() {
	rootCmd.AddCommand(addNewsboatCmd)

	addNewsboatCmd.Flags().BoolP(EditFlagName, EditFlagShort, false, "Edit the new bookmark in a text editor")
}
