package db

import (
	"errors"
	"fmt"
	"sort"
	"strings"

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
	Number      int     `json:"index"`
	Url         string  `json:"uri"`
	Name        string  `json:"title"`
	Description string  `json:"description"`
	Tags        TagList `json:"tags"`
}

func (bm *Bookmark) NameMatches(query string) bool {
	return subStringMatches(query, bm.Name)
}

func (bm *Bookmark) UrlMatches(query string) bool {
	return subStringMatches(query, bm.Url)
}

func (bm *Bookmark) DescriptionMatches(query string) bool {
	return subStringMatches(query, bm.Description)
}

func (bm *Bookmark) FormatAsInteractiveFileString() string {
	return fmt.Sprintf(
		BookmarkInteractiveFileFormatString,
		bm.Name, bm.Url, bm.Description, bm.Tags.MultilineString())
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
			bm.Url = l
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
	return bmks.Bookmarks[i].Number < bmks.Bookmarks[j].Number
}

func (bmks *BookmarkLibrary) Swap(i, j int) {
	bmks.Bookmarks[i], bmks.Bookmarks[j] = bmks.Bookmarks[j], bmks.Bookmarks[i]
}

func (bmks *BookmarkLibrary) Verify() error {
	// Get count of each bookmark number used
	counts := make(map[int]int)
	for _, bm := range bmks.Bookmarks {
		counts[bm.Number]++
	}

	for number, count := range counts {
		if count > 1 {
			return errors.New(fmt.Sprintf("Bookmark number %d used %d times", number, count))
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
		Number: maxNumber + 1,
		Name:   "Untitled",
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
	i := sort.Search(bmks.Len(), func(i int) bool {
		return bmks.Bookmarks[i].Number >= number
	})

	if i < bmks.Len() && bmks.Bookmarks[i].Number == number {
		return i, nil
	} else {
		return 0, errors.New(fmt.Sprintf("No bookmark with number %d found", number))
	}
}
