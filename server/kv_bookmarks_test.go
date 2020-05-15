package main

import (
	"testing"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getTestBookmarks() *Bookmarks {
	api := makeAPIMock()
	api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
	p := makePlugin(api)
	bmarks := NewBookmarksWithUser(p.API, UserID)

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

	// no title provided
	b3 := &Bookmark{
		PostID:     "ID3",
		CreateAt:   model.GetMillis(),
		ModifiedAt: model.GetMillis(),
	}

	_ = bmarks.add(b1)
	_ = bmarks.add(b2)
	_ = bmarks.add(b3)

	return bmarks
}

func TestBookmarks_get(t *testing.T) {
	bmarks := getTestBookmarks()
	assert.Equal(t, 3, len(bmarks.ByID))
	bmark := bmarks.get("ID3")
	assert.Equal(t, "", bmark.getTitle())
}

func TestBookmarks_add(t *testing.T) {
	b4 := &Bookmark{PostID: "ID4", Title: "Title4"}
	bmarks := getTestBookmarks()
	assert.Equal(t, 3, len(bmarks.ByID))
	err := bmarks.add(b4)
	assert.Nil(t, err)
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
	_, exists := bmarks.exists("ID2")
	assert.Equal(t, true, exists)
}

func TestBookmarks_updateTimes(t *testing.T) {
	bmarks := getTestBookmarks()

	// bmark has been initialized. times not yet added
	b1 := bmarks.get("ID1")
	assert.Equal(t, 0, int(b1.CreateAt))
	assert.Equal(t, 0, int(b1.ModifiedAt))

	// bmark has been added and times added
	bmarks.updateTimes("ID1")
	bmarks.get("ID1")
	assert.Greater(t, int(b1.ModifiedAt), 0)
	assert.Equal(t, int(b1.ModifiedAt), int(b1.CreateAt))

	// bmark was already saved and modified time updates
	time.Sleep(time.Millisecond)
	bmarks.updateTimes("ID2")
	b2 := bmarks.get("ID2")
	assert.Greater(t, b2.ModifiedAt, b2.CreateAt)
}
