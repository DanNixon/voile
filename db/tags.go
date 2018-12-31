package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type TagList struct {
	Tags []string
}

func (tl *TagList) Clear() {
	tl.Tags = []string{}
}

func (tl TagList) MarshalJSON() ([]byte, error) {
	return json.Marshal(tl.Tags)
}

func (tl *TagList) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &tl.Tags)
}

func (tl *TagList) Len() int {
	return len(tl.Tags)
}

func (tl *TagList) Less(i, j int) bool {
	return strings.Compare(tl.Tags[i], tl.Tags[j]) < 0
}

func (tl *TagList) Swap(i, j int) {
	tl.Tags[i], tl.Tags[j] = tl.Tags[j], tl.Tags[i]
}

func (tl TagList) String() string {
	return strings.Join(tl.Tags, ", ")
}

func (tl TagList) MultilineString() string {
	return strings.Join(tl.Tags, "\n")
}

func (tl *TagList) ContainsAllTags(tags []string) bool {
	if len(tags) == 0 {
		return false
	}

	for _, t := range tags {
		if _, err := tl.search(t); err != nil {
			return false
		}
	}

	return true
}

func (tl *TagList) Append(tag string) {
	tag = strings.TrimSpace(tag)
	if len(tag) > 0 {
		if _, err := tl.search(tag); err != nil {
			tl.Tags = append(tl.Tags, tag)
			sort.Sort(tl)
		}
	}
}

func (tl *TagList) AppendFromString(tagStr string) {
	for _, tag := range strings.Split(tagStr, ",") {
		tl.Append(tag)
	}
}

func (tl *TagList) Remove(tag string) {
	if len(tag) > 0 {
		if i, err := tl.search(tag); err == nil {
			tl.Tags = append(tl.Tags[:i], tl.Tags[i+1:]...)
		}
	}
}

func (tl *TagList) search(tag string) (int, error) {
	i := sort.Search(tl.Len(), func(i int) bool {
		return tl.Tags[i] >= tag
	})

	if i < tl.Len() && tl.Tags[i] == tag {
		return i, nil
	}

	return 0, errors.New(fmt.Sprintf("No tag found matching %s", tag))
}

type TagCount map[string]int

type AllTags struct {
	Tags  TagList
	Count TagCount
}
