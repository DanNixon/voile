package cmd

import (
	"bytes"
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
		// Get bookmark number
		bmks := ReadBookmarksFromFile()

		// Get action flags
		copyFlag, _ := cmd.Flags().GetBool(CopyFlagName)
		openFlag, _ := cmd.Flags().GetBool(OpenFlagName)

		// Buffer for clipboard string
		var clipboardBuffer bytes.Buffer

		if cmd.Flags().Changed(NumberFlagName) {
			// Get bookmark entry
			bookmarkNumber, _ := cmd.Flags().GetInt(NumberFlagName)
			bm, err := bmks.GetByNumber(bookmarkNumber)
			CheckError(err)

			// Print bookmark to console
			fmt.Println(FormatBookmark(bm, 0))

			// Buffer URL for clipboard copy
			if copyFlag {
				clipboardBuffer.WriteString(bm.Url)
			}

			// Open URL in browser
			if openFlag {
				open.Run(bm.Url)
			}
		} else {
			// Get filter options
			tags, _ := cmd.Flags().GetStringSlice(TagsFlagName)
			name, _ := cmd.Flags().GetString(NameFlagName)
			url, _ := cmd.Flags().GetString(UrlFlagName)
			desc, _ := cmd.Flags().GetString(DescFlagName)

			// Setup filtering
			d := FilteringOptions{
				[]FilterCase{
					FilterCase{
						cmd.Flags().Changed(TagsFlagName),
						func(b *db.Bookmark) bool { return b.Tags.ContainsAllTags(tags) },
					},
					FilterCase{
						cmd.Flags().Changed(NameFlagName),
						func(b *db.Bookmark) bool { return b.NameMatches(name) },
					},
					FilterCase{
						cmd.Flags().Changed(UrlFlagName),
						func(b *db.Bookmark) bool { return b.UrlMatches(url) },
					},
					FilterCase{
						cmd.Flags().Changed(DescFlagName),
						func(b *db.Bookmark) bool { return b.DescriptionMatches(desc) },
					},
				},
			}

			i := 0
			for _, bm := range bmks.Bookmarks {
				// Check if this bookmark should be excluded from the results
				if !includeBookmark(&d, &bm) {
					continue
				}

				// Print bookmark to console
				if i > 0 {
					fmt.Println()
				}
				fmt.Println(FormatBookmark(&bm, i))

				// Buffer URLs for clipboard copy
				if copyFlag {
					if i > 0 {
						clipboardBuffer.WriteString("\n")
					}
					clipboardBuffer.WriteString(bm.Url)
				}

				// Open URL in browser
				if openFlag {
					open.Run(bm.Url)
				}

				i++
			}
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
