package bookmarks

import (
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi/mock_pluginapi"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/assert"
)

const UserID = "UserID"

func getTestBookmarks() *Bookmarks {
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

	bmarks := NewBookmarks(UserID)
	bmarks.ByID[b1.PostID] = b1
	bmarks.ByID[b2.PostID] = b2
	bmarks.ByID[b3.PostID] = b3

	return bmarks
}

func TestBookmarks_get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPluginAPI := mock_pluginapi.NewMockAPI(ctrl)

	bmarks := getTestBookmarks()
	bmarks.api = mockPluginAPI

	assert.Equal(t, 3, len(bmarks.ByID))
	bmark, _ := bmarks.GetBookmark("ID3")
	assert.Equal(t, "", bmark.GetTitle())
}

func TestBookmarks_add(t *testing.T) {
	b4 := &Bookmark{PostID: "ID4", Title: "Title4"}
	bmarks := getTestBookmarks()
	assert.Equal(t, 3, len(bmarks.ByID))
	bmarks.ByID[b4.PostID] = b4
	assert.Equal(t, 4, len(bmarks.ByID))
}

func TestBookmarks_delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPluginAPI := mock_pluginapi.NewMockAPI(ctrl)

	mockPluginAPI.EXPECT().KVSet(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	bmarks := getTestBookmarks()
	bmarks.api = mockPluginAPI
	assert.Equal(t, 3, len(bmarks.ByID))
	_ = bmarks.DeleteBookmark("ID2")
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
	b1, _ := bmarks.GetBookmark("ID1")
	assert.Equal(t, 0, int(b1.CreateAt))
	assert.Equal(t, 0, int(b1.ModifiedAt))

	// bmark has been added and times added
	bmarks.updateTimes("ID1")
	assert.Greater(t, int(b1.ModifiedAt), 0)
	assert.Equal(t, int(b1.ModifiedAt), int(b1.CreateAt))

	// bmark was already saved and modified time updates
	time.Sleep(time.Millisecond)
	bmarks.updateTimes("ID2")
	b2, _ := bmarks.GetBookmark("ID2")
	assert.Greater(t, b2.ModifiedAt, b2.CreateAt)
}
