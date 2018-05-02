package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	. "github.com/logrusorgru/aurora"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/DanNixon/voile/db"

	"github.com/DanNixon/voile/tui"
)

const (
	BookmarksFileConfigEntry = "bookmark_file"
)

const (
	CopyFlagName  = "copy"
	CopyFlagShort = "c"

	OpenFlagName  = "open"
	OpenFlagShort = "o"

	EditFlagName  = "edit"
	EditFlagShort = "e"

	ForceFlagName = "force"

	NumberFlagName  = "number"
	NumberFlagShort = "n"

	TagsFlagName  = "tags"
	TagsFlagShort = "t"

	NameFlagName  = "name"
	NameFlagShort = "n"

	TitleNameFlagName  = "autoname"
	TitleNameFlagShort = "a"

	UrlFlagName  = "url"
	UrlFlagShort = "u"

	DescFlagName  = "desc"
	DescFlagShort = "d"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var bookmarkFormatStr = fmt.Sprintf(
	"%s %s %s%s%s\n  %s %s\n  %s %s",
	Black("%d."), Bold(Green("%s")), Black("["), Bold(Cyan("%d")), Black("]"),
	Red(">"), Brown("%s"),
	Red("#"), Blue("%s"))

func FormatBookmark(bm *db.Bookmark, index int) string {
	return fmt.Sprintf(
		bookmarkFormatStr,
		index, bm.Name, bm.Number, bm.Url, bm.Tags)
}

func EditBookmarkInEditor(bmks *db.BookmarkLibrary, bm *db.Bookmark) {
	var err error

	// Generate default bookmark string
	bmStr := bm.FormatAsInteractiveFileString()

	// Add existing tags commented at end of string
	bmStr += "# Existing tags:\n"
	for _, t := range bmks.GetAllTags().Tags.Tags {
		bmStr += "# " + t + "\n"
	}

	// Launch editor with existing bookmark data
	bmStr, err = tui.EditText(bmStr)
	CheckError(err)

	// Update bookmark with modifications
	err = bm.UpdateFromInteractiveFileString(bmStr)
	CheckError(err)
}

func IsValidBookmarkNumberArgument(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires exactly least one arg")
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	if index < 0 {
		return errors.New("Bookmark number must be >= 0")
	}

	return nil
}

func ReadBookmarksFromFile() db.BookmarkLibrary {
	var bmks db.BookmarkLibrary

	// Read entire JSON file to string
	raw, err := ioutil.ReadFile(viper.GetString(BookmarksFileConfigEntry))
	CheckError(err)

	// Load bookmarks from JSON
	json.Unmarshal(raw, &(bmks.Bookmarks))

	// Validate the loaded data
	err = bmks.Verify()
	CheckError(err)

	// Sort (by natural order of bookmark number)
	sort.Sort(&bmks)

	return bmks
}

func SaveBookmarksToFile(bmks *db.BookmarkLibrary) {
	// Write bookmarks to indented JSON string
	raw, err := json.MarshalIndent(bmks.Bookmarks, "", "  ")
	CheckError(err)

	// Write JSON string to file
	err = ioutil.WriteFile(
		viper.GetString(BookmarksFileConfigEntry), []byte(raw), 0644)
	CheckError(err)
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initBookmarksFile)
}

func initConfig() {
	// Get home directory
	home, err := homedir.Dir()
	CheckError(err)

	// Set default bookmarks file
	viper.SetDefault(
		BookmarksFileConfigEntry,
		filepath.Join(home, ".voile_bookmarks"))

	// Set config file location
	viper.SetConfigType("yaml")
	viper.SetConfigName(".voile")
	viper.AddConfigPath(home)
	viper.ReadInConfig()
}

func initBookmarksFile() {
	// Create an empty bookmarks file if one does not already exist
	if _, err := os.Stat(viper.GetString(BookmarksFileConfigEntry)); err != nil {
		var bmks db.BookmarkLibrary
		SaveBookmarksToFile(&bmks)
	}
}
