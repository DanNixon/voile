package db

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

const BookmarkInteractiveFileFormatString = `# Bookmark
# Headings must not be edited
# Lines starting with a # are ignored
## Title
%s
## URL
%s
## Description
%s
## Tags (comma and newline separated)
%s
`

const (
	BookmarkInteractiveFileNameHeader        = "## Title"
	BookmarkInteractiveFileUrlHeader         = "## URL"
	BookmarkInteractiveFileDescriptionHeader = "## Description"
	BookmarkInteractiveFileTagsHeader        = "## Tags"
)

func subStringMatches(query, source string) bool {
	m := search.New(language.Und, search.IgnoreCase)
	i, _ := m.IndexString(source, query)
	return i != -1
}

type Bookmark struct {
	Number      int       `json:"index"`
	Url         Url       `json:"uri"`
	Name        string    `json:"title"`
	Description string    `json:"description"`
	Tags        TagList   `json:"tags"`
	WhenAdded   time.Time `json:"whenAdded"`
}

func (bm *Bookmark) NameMatches(query string) bool {
	return subStringMatches(query, bm.Name)
}

func (bm *Bookmark) UrlMatches(query string) bool {
	return subStringMatches(query, bm.Url.String())
}

func (bm *Bookmark) DescriptionMatches(query string) bool {
	return subStringMatches(query, bm.Description)
}

func (bm *Bookmark) FormatAsInteractiveFileString() string {
	return fmt.Sprintf(
		BookmarkInteractiveFileFormatString,
		bm.Name, bm.Url.String(), bm.Description, bm.Tags.MultilineString())
}

func (bm *Bookmark) UpdateFromInteractiveFileString(data string) error {
	lines := strings.Split(data, "\n")
	header := ""

	var bookmarkInteractiveFileHeaderStrings = [...]string{
		BookmarkInteractiveFileNameHeader,
		BookmarkInteractiveFileUrlHeader,
		BookmarkInteractiveFileDescriptionHeader,
		BookmarkInteractiveFileTagsHeader,
	}

	bm.Description = ""

	for _, l := range lines {
		l = strings.TrimSpace(l)

		// Match header
		for _, h := range bookmarkInteractiveFileHeaderStrings {
			if strings.HasPrefix(l, h) {
				header = h
				break
			}
		}

		// Ignore commented out lines
		if strings.HasPrefix(l, "#") {
			continue
		}

		// Parse the line based on which valid line it was
		switch header {
		case BookmarkInteractiveFileNameHeader:
			bm.Name = l
		case BookmarkInteractiveFileUrlHeader:
			err := bm.Url.Parse(l)
			if err != nil {
				return err
			}
		case BookmarkInteractiveFileDescriptionHeader:
			if len(l) > 0 {
				bm.Description += l
			} else {
				bm.Description += "\n"
			}
		case BookmarkInteractiveFileTagsHeader:
			bm.Tags.AppendFromString(l)
		}
	}

	return nil
}

type BookmarkLibrary struct {
	Bookmarks []Bookmark
}

func (bmks *BookmarkLibrary) Len() int {
	return len(bmks.Bookmarks)
}

func (bmks *BookmarkLibrary) Less(i, j int) bool {
	return bmks.Bookmarks[i].WhenAdded.Before(bmks.Bookmarks[j].WhenAdded)
}

func (bmks *BookmarkLibrary) Swap(i, j int) {
	bmks.Bookmarks[i], bmks.Bookmarks[j] = bmks.Bookmarks[j], bmks.Bookmarks[i]
}

func (bmks *BookmarkLibrary) Verify() error {
	// Get count of each bookmark number and URL used
	numberCounts := make(map[int]int)
	urlCounts := make(map[string]int)
	for _, bm := range bmks.Bookmarks {
		numberCounts[bm.Number]++
		urlCounts[bm.Url.String()]++
	}

	for number, count := range numberCounts {
		if count > 1 {
			return errors.New(fmt.Sprintf("Bookmark number %d used %d times", number, count))
		}
	}

	for url, count := range urlCounts {
		if count > 1 {
			return errors.New(fmt.Sprintf("Bookmark URL %s used %d times", url, count))
		}
	}

	return nil
}

func (bmks *BookmarkLibrary) GetByNumber(number int) (*Bookmark, error) {
	i, err := bmks.searchByNumber(number)
	if err == nil {
		return &(bmks.Bookmarks[i]), nil
	} else {
		return nil, err
	}
}

func (bmks *BookmarkLibrary) DeleteByNumber(number int) error {
	i, err := bmks.searchByNumber(number)
	if err != nil {
		return err
	}

	bmks.Bookmarks = append(bmks.Bookmarks[:i], bmks.Bookmarks[i+1:]...)
	return nil
}

func (bmks *BookmarkLibrary) NewEntry() *Bookmark {
	maxNumber := 0
	for _, bm := range bmks.Bookmarks {
		if bm.Number > maxNumber {
			maxNumber = bm.Number
		}
	}

	bm := Bookmark{
		Number:    maxNumber + 1,
		Name:      "Untitled",
		WhenAdded: time.Now(),
	}
	bmks.Bookmarks = append(bmks.Bookmarks, bm)

	return &(bmks.Bookmarks[bmks.Len()-1])
}

func (bmks *BookmarkLibrary) GetAllTags() AllTags {
	var tags AllTags
	tags.Count = make(TagCount)

	for _, bm := range bmks.Bookmarks {
		for _, t := range bm.Tags.Tags {
			tags.Tags.Append(t)
			tags.Count[t]++
		}
	}

	sort.Sort(&tags.Tags)

	return tags
}

func (bmks *BookmarkLibrary) searchByNumber(number int) (int, error) {
	for idx, bm := range bmks.Bookmarks {
		if bm.Number == number {
			return idx, nil
		}
	}

	return 0, errors.New(fmt.Sprintf("No bookmark with number %d found", number))
}
