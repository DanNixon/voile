package db_test

import (
	"errors"
	"net/url"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/DanNixon/voile/db"
)

const TestBookmarkInteractiveFileExpectedString = `# Bookmark
# Headings must not be edited
# Lines starting with a # are ignored
## Title
BBC
## URL
https://bbc.co.uk
## Description
BBC News and Weather
## Tags (comma and newline separated)
news
weather
`

const TestBookmarkInteractiveFileExpectedStringMultiline = `# Bookmark
# Headings must not be edited
# Lines starting with a # are ignored
## Title
BBC
## URL
https://bbc.co.uk
## Description
BBC News and Weather.
(and sport)
## Tags (comma and newline separated)
news,weather
sport
`

var testBookmark = db.Bookmark{
	Number: 0,
	Name:   "BBC",
	Url: db.Url{
		url.URL{
			Scheme: "https",
			Host:   "bbc.co.uk",
		}},
	Description: "BBC News and Weather",
	Tags: db.TagList{
		Tags: []string{"news", "weather"},
	},
}

func createTestLibrary() db.BookmarkLibrary {
	return db.BookmarkLibrary{
		Bookmarks: []db.Bookmark{
			{
				Number: 1,
				Name:   "one",
				Url: db.Url{
					url.URL{
						Scheme: "https",
						Host:   "github.com",
					}},
				Tags: db.TagList{
					Tags: []string{"news", "weather"},
				},
			},
			{
				Number: 2,
				Name:   "two",
				Url: db.Url{
					url.URL{
						Scheme: "https",
						Host:   "facebook.com",
					}},
				Tags: db.TagList{
					Tags: []string{"software"},
				},
			},
			{
				Number: 3,
				Name:   "three",
				Url: db.Url{
					url.URL{
						Scheme: "https",
						Host:   "bbc.co.uk",
					}},
				Tags: db.TagList{
					Tags: []string{"news", "software"},
				},
			},
		},
	}
}

func TestBookmarkInit(t *testing.T) {
	var bm db.Bookmark
	assert.Equal(t, 0, bm.Number)
	assert.Equal(t, "", bm.Name)
	assert.Equal(t, "", bm.Url.String())
	assert.Equal(t, "", bm.Description)
	assert.Equal(t, db.TagList{}, bm.Tags)
}

func TestBookmarkHasName(t *testing.T) {
	var bm db.Bookmark
	assert.False(t, bm.HasName())
	bm.Name = "Test"
	assert.True(t, bm.HasName())
}

func TestBookmarkParseUrl(t *testing.T) {
	var bm db.Bookmark
	assert.Nil(t, bm.Url.Parse("https://github.com"))
	assert.Equal(t, "https://github.com", bm.Url.String())
}

func TestBookmarkNameMatches(t *testing.T) {
	assert.False(t, testBookmark.NameMatches(""))
	assert.False(t, testBookmark.NameMatches("ITV"))

	assert.True(t, testBookmark.NameMatches("B"))
	assert.True(t, testBookmark.NameMatches("c"))
	assert.True(t, testBookmark.NameMatches("Bb"))
	assert.True(t, testBookmark.NameMatches("bC"))
	assert.True(t, testBookmark.NameMatches("bbc"))
}

func TestBookmarkUrlMatches(t *testing.T) {
	assert.False(t, testBookmark.DescriptionMatches(""))
	assert.False(t, testBookmark.DescriptionMatches("sport"))

	assert.True(t, testBookmark.DescriptionMatches("news and weather"))
}

func TestBookmarkDescriptionMatches(t *testing.T) {
	assert.False(t, testBookmark.UrlMatches(""))
	assert.False(t, testBookmark.UrlMatches("ITV"))

	assert.True(t, testBookmark.UrlMatches("bbc.co.uk"))
	assert.True(t, testBookmark.UrlMatches(".co.uk"))
	assert.True(t, testBookmark.UrlMatches("https"))
}

func TestBookmarkFormatAsInteractiveFileString(t *testing.T) {
	testStr := testBookmark.FormatAsInteractiveFileString()
	assert.Equal(t, TestBookmarkInteractiveFileExpectedString, testStr)
}

func TestBookmarkUpdateFromInteractiveFileString(t *testing.T) {
	var bm db.Bookmark
	bm.UpdateFromInteractiveFileString(TestBookmarkInteractiveFileExpectedString)
	assert.Equal(t, testBookmark, bm)
}

func TestBookmarkUpdateFromInteractiveFileStringMultiline(t *testing.T) {
	var bm db.Bookmark
	bm.UpdateFromInteractiveFileString(TestBookmarkInteractiveFileExpectedStringMultiline)

	var testBookmark = db.Bookmark{
		Number: 0,
		Name:   "BBC",
		Url: db.Url{url.URL{
			Scheme: "https",
			Host:   "bbc.co.uk",
		}},
		Description: "BBC News and Weather.\n(and sport)",
		Tags: db.TagList{
			Tags: []string{"news", "sport", "weather"},
		},
	}
	assert.Equal(t, testBookmark, bm)
}

func TestBookmarkLibraryInit(t *testing.T) {
	var bmks db.BookmarkLibrary
	assert.Equal(t, 0, bmks.Len())
}

func TestBookmarkLibrarySort(t *testing.T) {
	bmks := db.BookmarkLibrary{
		Bookmarks: []db.Bookmark{
			{Number: 7, WhenAdded: time.Date(2018, time.November, 2, 10, 0, 0, 0, time.UTC)},
			{Number: 9, WhenAdded: time.Date(2018, time.December, 1, 10, 0, 0, 0, time.UTC)},
			{Number: 0, WhenAdded: time.Date(2018, time.July, 1, 10, 0, 0, 0, time.UTC)},
			{Number: 4, WhenAdded: time.Date(2018, time.November, 1, 10, 0, 0, 0, time.UTC)},
			{Number: 3, WhenAdded: time.Date(2018, time.October, 1, 10, 0, 0, 0, time.UTC)},
		},
	}

	sort.Sort(&bmks)

	assert.Equal(t, 0, bmks.Bookmarks[0].Number)
	assert.Equal(t, 3, bmks.Bookmarks[1].Number)
	assert.Equal(t, 4, bmks.Bookmarks[2].Number)
	assert.Equal(t, 7, bmks.Bookmarks[3].Number)
	assert.Equal(t, 9, bmks.Bookmarks[4].Number)
}

func TestBookmarkLibraryVerify(t *testing.T) {
	bmks := createTestLibrary()

	assert.Nil(t, bmks.Verify())
}

func TestBookmarkLibraryVerifyDuplicateNumbers(t *testing.T) {
	bmks := db.BookmarkLibrary{
		Bookmarks: []db.Bookmark{
			{Number: 7},
			{Number: 9},
			{Number: 7},
		},
	}

	err := bmks.Verify()

	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Bookmark number 7 used 2 times"), err)
}

func TestBookmarkLibraryGetByNumber(t *testing.T) {
	bmks := createTestLibrary()

	bm, err := bmks.GetByNumber(2)
	assert.Nil(t, err)
	assert.Equal(t, "two", bm.Name)
}

func TestBookmarkLibraryGetByNumberInvalid(t *testing.T) {
	bmks := createTestLibrary()

	bm, err := bmks.GetByNumber(7)
	assert.NotNil(t, err)
	assert.Nil(t, bm)
}

func TestBookmarkLibraryDeleteByNumber(t *testing.T) {
	bmks := createTestLibrary()
	oldLen := bmks.Len()

	bmks.DeleteByNumber(3)
	assert.Equal(t, oldLen-1, bmks.Len())
}

func TestBookmarkLibraryDeleteByNumberInvalid(t *testing.T) {
	bmks := createTestLibrary()
	oldLen := bmks.Len()

	bmks.DeleteByNumber(7)
	assert.Equal(t, oldLen, bmks.Len())
}

func TestBookmarkLibraryNewEntry(t *testing.T) {
	bmks := createTestLibrary()
	oldLen := bmks.Len()

	newBm := bmks.NewEntry()
	newBm.Name = "New Bookmark"
	assert.Equal(t, oldLen+1, bmks.Len())

	bm, err := bmks.GetByNumber(newBm.Number)
	assert.Nil(t, err)
	assert.Equal(t, "New Bookmark", bm.Name)
}

func TestBookmarkLibraryGetAllTags(t *testing.T) {
	bmks := createTestLibrary()

	tags := bmks.GetAllTags()

	assert.Equal(t, []string{"news", "software", "weather"}, tags.Tags.Tags)
	assert.Equal(t, 2, tags.Count["news"])
	assert.Equal(t, 1, tags.Count["weather"])
	assert.Equal(t, 2, tags.Count["software"])
}
