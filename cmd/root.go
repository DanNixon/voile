package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"

	"github.com/DanNixon/voile/db"
)

type FilteringOptions struct {
	Cases []FilterCase
}

type FilterCase struct {
	UserCares bool
	Func      func(*db.Bookmark) bool
}

var rootCmd = &cobra.Command{
	Use:   "voile",
	Short: "Query bookmark library",
	Long:  `Query bookmark library and open results.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load bookmarks from file
		bmks := ReadBookmarksFromFile()

		// Get action flags
		copyFlag, _ := cmd.Flags().GetBool(CopyFlagName)
		openFlag, _ := cmd.Flags().GetBool(OpenFlagName)

		// Get display flags
		jsonFlag, _ := cmd.Flags().GetBool(JsonFlagName)

		// Get filter options
		bookmarkNumber, _ := cmd.Flags().GetInt(NumberFlagName)
		tags, _ := cmd.Flags().GetStringSlice(TagsFlagName)
		name, _ := cmd.Flags().GetString(NameFlagName)
		url, _ := cmd.Flags().GetString(UrlFlagName)
		desc, _ := cmd.Flags().GetString(DescFlagName)

		// Setup filtering
		d := FilteringOptions{
			[]FilterCase{
				{
					cmd.Flags().Changed(NumberFlagName),
					func(b *db.Bookmark) bool { return b.Number == bookmarkNumber },
				},
				{
					cmd.Flags().Changed(TagsFlagName),
					func(b *db.Bookmark) bool { return b.Tags.ContainsAllTags(tags) },
				},
				{
					cmd.Flags().Changed(NameFlagName),
					func(b *db.Bookmark) bool { return b.NameMatches(name) },
				},
				{
					cmd.Flags().Changed(UrlFlagName),
					func(b *db.Bookmark) bool { return b.UrlMatches(url) },
				},
				{
					cmd.Flags().Changed(DescFlagName),
					func(b *db.Bookmark) bool { return b.DescriptionMatches(desc) },
				},
			},
		}

		// Buffer for clipboard string
		var clipboardBuffer bytes.Buffer

		// Collect bookmarks for JSON output
		var filteredBookmarks []db.Bookmark

		// Filter bookmarks
		i := 0
		for _, bm := range bmks.Bookmarks {
			// Check if this bookmark should be excluded from the results
			if !includeBookmark(&d, &bm) {
				continue
			}

			if jsonFlag {
				// Add bookmark to filtered list for JSON output
				filteredBookmarks = append(filteredBookmarks, bm)
			} else {
				// Output to console in standard format
				if i > 0 {
					fmt.Println()
				}
				fmt.Println(FormatBookmark(&bm, i))
			}

			// Buffer URLs for clipboard copy
			if copyFlag {
				if i > 0 {
					clipboardBuffer.WriteString("\n")
				}
				clipboardBuffer.WriteString(bm.Url.String())
			}

			// Open URL in browser
			if openFlag {
				open.Run(bm.Url.String())
			}

			i++
		}

		// Print bookmark to console
		if jsonFlag {
			raw, err := json.MarshalIndent(filteredBookmarks, "", "  ")
			CheckError(err)
			fmt.Printf("%s\n", raw)
		}

		// Write URLs to clipboard
		if copyFlag {
			clipboard.WriteAll(clipboardBuffer.String())
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	CheckError(err)
}

func init() {
	rootCmd.Flags().IntP(NumberFlagName, NumberFlagShort, 0, "Get bookmark by number")
	rootCmd.Flags().StringSliceP(TagsFlagName, TagsFlagShort, []string{}, "Get bookmarks by tags")
	rootCmd.Flags().StringP(NameFlagName, "s", "", "Search in name")
	rootCmd.Flags().StringP(UrlFlagName, UrlFlagShort, "", "Search in URL")
	rootCmd.Flags().StringP(DescFlagName, DescFlagShort, "", "Search in description")

	rootCmd.Flags().BoolP(OpenFlagName, OpenFlagShort, false, "Open bookmarks in browser")
	rootCmd.Flags().BoolP(CopyFlagName, CopyFlagShort, false, "Copy bookmark URLs to clipboard")

	rootCmd.Flags().BoolP(JsonFlagName, JsonFlagShort, false, "Output in JSON format")
}

func includeBookmark(query *FilteringOptions, bm *db.Bookmark) bool {
	res := true
	for _, q := range query.Cases {
		if q.UserCares {
			res = res && q.Func(bm)
		}
	}
	return res
}
