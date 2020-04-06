package main

import (
	"testing"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/assert"
)

func getTestBookmarks() *Bookmarks {
	bmarks := NewBookmarks()

	b1 := &Bookmark{
		PostID: "ID1",
		Title:  "Title1 - New Bookmark - times are zero",
	}
	b2 := &Bookmark{
		PostID:     "ID2",
		Title:      "Title2 - bookmarks initialized. Times created and same",
		CreateAt:   model.GetMillis(),
		ModifiedAt: model.GetMillis(),
	}

	b3 := &Bookmark{
		PostID:     "ID3",
		Title:      "Title3 - bookmarks already updated once",
		CreateAt:   model.GetMillis(),
		ModifiedAt: model.GetMillis(),
	}

	bmarks.add(b1)
	bmarks.add(b2)
	bmarks.add(b3)

	return bmarks
}

func TestBookmarks_get(t *testing.T) {
	bmarks := getTestBookmarks()
	assert.Equal(t, 3, len(bmarks.ByID))
	bmark := bmarks.get("ID3")
	assert.Equal(t, "Title3 - bookmarks already updated once", bmark.Title)
}

func TestBookmarks_add(t *testing.T) {
	b4 := &Bookmark{PostID: "ID4", Title: "Title4"}
	bmarks := getTestBookmarks()
	assert.Equal(t, 3, len(bmarks.ByID))
	bmarks.add(b4)
	assert.Equal(t, 4, len(bmarks.ByID))
}

func TestBookmarks_delete(t *testing.T) {
	bmarks := getTestBookmarks()
	assert.Equal(t, 3, len(bmarks.ByID))
	bmarks.delete("ID2")
	assert.Equal(t, 2, len(bmarks.ByID))
}

func TestBookmarks_exists(t *testing.T) {
	bmarks := getTestBookmarks()
	exists := bmarks.exists("ID2")
	assert.Equal(t, true, exists)
}

func TestBookmarks_updateTimes(t *testing.T) {
	bmarks := getTestBookmarks()

	// bmark has been initialized. times not yet added
	b1 := bmarks.get("ID1")
	assert.Equal(t, 0, int(b1.CreateAt))
	assert.Equal(t, 0, int(b1.ModifiedAt))

	// bmark has been added and times added
	b1 = bmarks.updateTimes("ID1")
	b1 = bmarks.get("ID1")
	assert.Greater(t, int(b1.ModifiedAt), 0)
	assert.Equal(t, int(b1.ModifiedAt), int(b1.CreateAt))

	// bmark was already saved and modified time updates
	b2 := bmarks.get("ID2")
	time.Sleep(time.Millisecond)
	b2 = bmarks.updateTimes("ID2")
	b2 = bmarks.get("ID2")
	assert.Greater(t, b2.ModifiedAt, b2.CreateAt)
}
