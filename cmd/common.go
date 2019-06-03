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
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-gitconfig"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

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

	JsonFlagName  = "json"
	JsonFlagShort = "j"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func FormatBookmark(bm *db.Bookmark, index int) string {
	// Generate name
	var nameStr aurora.Value
	if bm.HasName() {
		nameStr = aurora.Green(bm.Name)
	} else {
		nameStr = aurora.Red("[untitled]")
	}

	// Name and URL
	retVal := fmt.Sprintf(
		"%s. %s [%s]\n  %s %s",
		strconv.Itoa(index), aurora.Bold(nameStr),
		aurora.Bold(aurora.Cyan(strconv.Itoa(bm.Number))),
		aurora.Red(">"), aurora.Brown(bm.Url.String()))

	// Tags (if set)
	if bm.Tags.Len() > 0 {
		retVal += fmt.Sprintf("\n  %s %s", aurora.Red("#"), aurora.Blue(bm.Tags))
	}

	// Description (if set)
	if len(bm.Description) > 0 {
		r := strings.NewReplacer("\n", "\n    ")
		retVal += fmt.Sprintf("\n  %s %s", aurora.Red("?"),
			r.Replace(bm.Description))
	}

	// Added timestamp
	retVal += fmt.Sprintf("\n  %s %s", aurora.Red("+"),
		aurora.Cyan(bm.WhenAdded.Format(time.UnixDate)))

	return retVal
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

	sort.Sort(&bmks)

	return bmks
}

func SaveBookmarksToFile(bmks *db.BookmarkLibrary) {
	// Validate the bookmarks before saving
	err := bmks.Verify()
	CheckError(err)

	// Write bookmarks to indented JSON string
	raw, err := json.MarshalIndent(bmks.Bookmarks, "", "  ")
	CheckError(err)

	// Write JSON string to file
	err = ioutil.WriteFile(
		viper.GetString(BookmarksFileConfigEntry), []byte(raw), 0644)
	CheckError(err)

	// Git commit
	err = CommitChangesToBookmarkFile()
	CheckError(err)
}

func GetBookmarksFileParentDirectory() string {
	filename := viper.GetString(BookmarksFileConfigEntry)
	return filepath.Dir(filename)
}

func IsBookmarksFileInGitRepository() bool {
	_, err := git.PlainOpen(GetBookmarksFileParentDirectory())
	return err == nil
}

func CommitChangesToBookmarkFile() error {
	gitDir := GetBookmarksFileParentDirectory()

	repo, err := git.PlainOpen(gitDir)
	if err != nil {
		return nil
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	file, err := filepath.Rel(gitDir, viper.GetString(BookmarksFileConfigEntry))
	if err != nil {
		return err
	}

	_, err = wt.Add(file)
	if err != nil {
		return err
	}

	username, _ := gitconfig.Username()
	email, _ := gitconfig.Email()

	_, err = wt.Commit("voile auto commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  username,
			Email: email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initBookmarksFile)
}

func initConfig() {
	// Setup environment variable config options
	viper.SetEnvPrefix("voile")
	viper.BindEnv(BookmarksFileConfigEntry)

	// Set default bookmarks file
	viper.SetDefault(BookmarksFileConfigEntry, "bookmarks.json")
}

func initBookmarksFile() {
	// Create an empty bookmarks file if one does not already exist
	if _, err := os.Stat(viper.GetString(BookmarksFileConfigEntry)); err != nil {
		var bmks db.BookmarkLibrary
		SaveBookmarksToFile(&bmks)
	}
}
