package cmd

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add URL",
	Short: "Add a new bookmark",
	Long:  `Adds a new bookmark, either by a set of flags or specifying fields via a text editor.`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get URL
		var url string
		copyFlag, _ := cmd.Flags().GetBool(CopyFlagName)
		if len(args) == 0 && copyFlag {
			// Copy URL from clipboard
			url, _ = clipboard.ReadAll()
		} else if len(args) == 1 && !copyFlag {
			// Get URL from argument
			url = args[0]
		} else {
			fmt.Println("Ambiguous URL source")
			os.Exit(1)
		}

		// Load bookmarks from file
		bmks := ReadBookmarksFromFile()

		// Create new bookmark entry
		bm := bmks.NewEntry()

		// Set URL
		bm.Url = url

		// Set name
		if cmd.Flags().Changed(NameFlagName) {
			bm.Name, _ = cmd.Flags().GetString(NameFlagName)
		}

		// Set description
		if cmd.Flags().Changed(DescFlagName) {
			bm.Description, _ = cmd.Flags().GetString(DescFlagName)
		}

		// Set tags
		if cmd.Flags().Changed(TagsFlagName) {
			tags, _ := cmd.Flags().GetStringSlice(TagsFlagName)
			for _, t := range tags {
				bm.Tags.Append(t)
			}
		}

		// Edit in editor if requested
		editFlag, _ := cmd.Flags().GetBool(EditFlagName)
		if editFlag {
			EditBookmarkInEditor(bm)
		}

		// Print bookmark to console
		fmt.Println(FormatBookmark(bm, 0))

		// Save bookmarks back to file
		SaveBookmarksToFile(&bmks)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP(CopyFlagName, CopyFlagShort, false, "Copy URL from clipboard")
	addCmd.Flags().BoolP(EditFlagName, EditFlagShort, false, "Add/edit the new bookmark in a text editor")

	addCmd.Flags().StringSliceP(TagsFlagName, TagsFlagShort, []string{}, "Tags")
	addCmd.Flags().StringP(NameFlagName, NameFlagShort, "", "Name")
	addCmd.Flags().StringP(DescFlagName, DescFlagShort, "", "Description")
}
