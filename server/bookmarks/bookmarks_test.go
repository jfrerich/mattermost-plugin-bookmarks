package bookmarks

import (
	"encoding/json"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi/mock_pluginapi"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"

	"github.com/stretchr/testify/assert"
)

//nolint
func makePlugin(api *plugintest.API) *plugin.MattermostPlugin {
	p := &plugin.MattermostPlugin{}
	p.SetAPI(api)
	return p
}

func TestStoreBookmarks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPluginAPI := mock_pluginapi.NewMockAPI(ctrl)

	// initialize test Bookmarks
	u1 := "userID1"

	b1 := &Bookmark{PostID: "ID1", Title: "Title1"}
	b2 := &Bookmark{PostID: "ID2", Title: "Title2"}

	// Add Bookmarks
	bmarks := NewBookmarks(u1)
	bmarks.ByID[b1.PostID] = b1
	bmarks.ByID[b2.PostID] = b2
	bmarks.api = mockPluginAPI

	// Markshal the bmarks and mock api call
	jsonBookmarks, err := json.Marshal(bmarks)
	assert.Nil(t, err)
	mockPluginAPI.EXPECT().KVSet("bookmarks_userID1", jsonBookmarks).Return(nil)

	// store bmarks using API
	err = bmarks.StoreBookmarks()
	assert.Nil(t, err)
}

func TestAddBookmark(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPluginAPI := mock_pluginapi.NewMockAPI(ctrl)

	// create some test bookmarks
	b1 := &Bookmark{PostID: "ID1", Title: "Title1"}
	b2 := &Bookmark{
		PostID:   "ID2",
		Title:    "Title2",
		LabelIDs: []string{"UUID1", "UUID2"},
	}
	b3 := &Bookmark{PostID: "ID3", Title: "Title3"}

	// User 1 has no bookmarks
	u1 := "userID1"
	u2 := "userID2"
	bmarksU1 := NewBookmarks(u1)
	bmarksU1.api = mockPluginAPI

	bmarksU2 := NewBookmarks(u2)
	bmarksU2.ByID[b1.PostID] = b1
	bmarksU2.ByID[b2.PostID] = b2
	bmarksU2.api = mockPluginAPI

	tests := []struct {
		name    string
		userID  string
		bmarks  *Bookmarks
		want    int
		wantErr bool
	}{
		{
			name:    "u1 no previous bookmarks  add one bookmark",
			userID:  u1,
			bmarks:  bmarksU1,
			wantErr: true,
			want:    1,
		},
		{
			name:    "u2 two previous bookmarks  add one bookmark",
			userID:  u2,
			bmarks:  bmarksU2,
			wantErr: true,
			want:    3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBmarks, err := json.Marshal(tt.bmarks)
			assert.Nil(t, err)

			mockPluginAPI.EXPECT().KVGet(GetBookmarksKey(UserID)).Return(jsonBmarks, nil).AnyTimes()
			mockPluginAPI.EXPECT().KVSet(gomock.Any(), gomock.Any()).Return(nil)

			// store bmarks using API
			err = tt.bmarks.AddBookmark(b3)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, len(tt.bmarks.ByID))
		})
	}
}

func TestDeleteBookmark(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPluginAPI := mock_pluginapi.NewMockAPI(ctrl)

	// create some test bookmarks
	b1 := &Bookmark{PostID: "ID1", Title: "Title1"}
	b2 := &Bookmark{
		PostID:   "ID2",
		Title:    "Title2",
		LabelIDs: []string{"UUID1", "UUID2"},
	}

	// User 1 has no bookmarks
	u1 := "userID1"
	u2 := "userID2"

	// User 2 has 2 existing bookmarks
	bmarksU1 := NewBookmarks(u1)
	bmarksU1.ByID[b1.PostID] = b1
	bmarksU1.api = mockPluginAPI

	bmarksU2 := NewBookmarks(u2)
	bmarksU2.ByID[b2.PostID] = b2
	bmarksU2.api = mockPluginAPI

	tests := []struct {
		name       string
		userID     string
		bmarks     *Bookmarks
		wantErrMsg string
		wantErr    bool
	}{
		{
			name:       "u1 no previous bookmarks  Error out",
			userID:     u1,
			bmarks:     bmarksU1,
			wantErr:    true,
			wantErrMsg: "Bookmark `ID2` does not exist",
		},
		{
			name:    "u2 two previous bookmarks  delete one bookmark",
			userID:  u2,
			bmarks:  bmarksU2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBmarks, err := json.Marshal(tt.bmarks)
			assert.Nil(t, err)

			if !tt.wantErr {
				mockPluginAPI.EXPECT().KVGet(GetBookmarksKey(tt.userID)).Return(jsonBmarks, nil).AnyTimes()
			}

			// not testing store in this test.  mock to accept anything
			mockPluginAPI.EXPECT().KVSet(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			err = tt.bmarks.DeleteBookmark(b2.PostID)
			if tt.wantErr {
				assert.Equal(t, err.Error(), tt.wantErrMsg)
				return
			}

			assert.Nil(t, err)
		})
	}
}
