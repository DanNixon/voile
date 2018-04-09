package db_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DanNixon/voile/db"
)

func TestTagListInit(t *testing.T) {
	var tl db.TagList
	assert.Equal(t, tl.Len(), 0)
	assert.Nil(t, tl.Tags)
}

func TestTagListSort(t *testing.T) {
	tl := db.TagList{
		Tags: []string{"bbb", "zzz", "iii"},
	}

	sort.Sort(&tl)
	assert.Equal(t, []string{"bbb", "iii", "zzz"}, tl.Tags)
}

func TestTagListAppend(t *testing.T) {
	var tl db.TagList

	tl.Append("one")
	assert.Equal(t, []string{"one"}, tl.Tags)

	tl.Append("two")
	assert.Equal(t, []string{"one", "two"}, tl.Tags)
}

func TestTagListAppendEmptyString(t *testing.T) {
	var tl db.TagList

	tl.Append("one")
	assert.Equal(t, []string{"one"}, tl.Tags)

	tl.Append("")
	assert.Equal(t, []string{"one"}, tl.Tags)
}

func TestTagListAppendDuplicate(t *testing.T) {
	var tl db.TagList

	tl.Append("one")
	assert.Equal(t, []string{"one"}, tl.Tags)

	tl.Append("one")
	assert.Equal(t, []string{"one"}, tl.Tags)
}

func TestTagListAppendDuplicateTrim(t *testing.T) {
	var tl db.TagList

	tl.Append(" one")
	assert.Equal(t, []string{"one"}, tl.Tags)

	tl.Append("one ")
	assert.Equal(t, []string{"one"}, tl.Tags)
}

func TestTagListRemove(t *testing.T) {
	tl := db.TagList{
		Tags: []string{"one", "two", "three"},
	}

	tl.Remove("two")
	assert.Equal(t, []string{"one", "three"}, tl.Tags)

	tl.Remove("one")
	assert.Equal(t, []string{"three"}, tl.Tags)
}

func TestTagListRemoveEmptyString(t *testing.T) {
	tl := db.TagList{
		Tags: []string{"one", "two", "three"},
	}

	tl.Remove("")
	assert.Equal(t, []string{"one", "two", "three"}, tl.Tags)
}

func TestTagListRemoveNotInList(t *testing.T) {
	tl := db.TagList{
		Tags: []string{"one", "two", "three"},
	}

	tl.Remove("eight")
	assert.Equal(t, []string{"one", "two", "three"}, tl.Tags)
}

func TestTagListContainsAllTags(t *testing.T) {
	var tl db.TagList
	tl.Append("one")
	tl.Append("two")
	tl.Append("three")

	assert.True(t, tl.ContainsAllTags([]string{"one"}))
	assert.True(t, tl.ContainsAllTags([]string{"one", "two"}))
	assert.False(t, tl.ContainsAllTags([]string{"one", "two", "nine"}))
}

func TestTagListContainsAllTagsEmpty(t *testing.T) {
	var tl db.TagList
	tl.Append("one")
	tl.Append("two")
	tl.Append("three")

	assert.False(t, tl.ContainsAllTags([]string{}))
}

func TestTagListContainsAllTagsEmptyString(t *testing.T) {
	var tl db.TagList
	tl.Append("one")
	tl.Append("two")
	tl.Append("three")

	assert.False(t, tl.ContainsAllTags([]string{""}))
}
