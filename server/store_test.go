package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// pluginWithMockedSubs returns mocked plugin for given subscriptions
func pluginWithMockedBookmarks(userID string, bookmark *Bookmark) *Plugin {
	p := &Plugin{}
	mockPluginAPI := &plugintest.API{}

	bmarks := Bookmarks{ByID: map[string]*Bookmark{}}
	bmarks.ByID[userID] = bookmark
	// bmarksD, _ := json.MarshalIndent(bmarks, "", "    ")
	// fmt.Printf("bmarks = %+v\n", string(bmarksD))

	// jsn, _ := json.Marshal(bmarks)
	mockPluginAPI.On("KVGet", getBookmarksKey(userID)).Return(nil, nil)
	mockPluginAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)

	// mockPluginAPI.On("KVSet", getBookmarksKey(userID)).Return(nil, nil)
	p.SetAPI(mockPluginAPI)

	return p
}

func TestPlugin_addToBookmarksForUser(t *testing.T) {
	type args struct {
		userID      string
		newBookmark *Bookmark
	}
	tests := []struct {
		name    string
		plugin  *Plugin
		args    args
		want    *Bookmark
		wantErr bool
	}{
		{
			name: "One Bookmark added",
			args: args{
				userID: "userID1",
				newBookmark: &Bookmark{
					PostID: "PostID_1",
					Title:  "Title_1",
				},
			},
			wantErr: true,
			plugin: pluginWithMockedBookmarks("userID1",
				&Bookmark{
					PostID: "PostID_1",
					Title:  "Title_1",
				},
			),
			want: &Bookmark{
				PostID: "PlkjostID_1",
				Title:  "Title_1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := tt.plugin.addToBookmarksForUser(tt.args.userID, tt.args.newBookmark); (err != nil) != tt.wantErr {
				gotD, _ := json.MarshalIndent(got, "", "    ")
				fmt.Printf("got = %+v\n", string(gotD))
				// bmarks, err := tt.plugin.getBookmarksForUser(tt.args.userID)
				// fmt.Printf("bmarks = %+v\n", bmarks)
				t.Errorf("Plugin.addToBookmarksForUser() error = %v, wantErr %v", err, tt.wantErr)
				assert.Equal(t, tt.want, got, "they should be the same")
				// fmt.Printf("bmarks = %+v\n", bmarks)
			}
		})
	}
}
