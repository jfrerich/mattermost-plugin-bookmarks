package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
)

// pluginWithMockedSubs returns mocked plugin for given subscriptions
func pluginWithMockedBookmarks(userID string, bookmarks []*Bookmark) *Plugin {
	p := &Plugin{}
	mockPluginAPI := &plugintest.API{}

	bmarks := Bookmarks{}
	bmarks.Bookmarks = bookmarks

	bmarksD, _ := json.MarshalIndent(bmarks, "", "    ")
	fmt.Printf("bmarks = %+v\n", string(bmarksD))

	jsn, _ := json.Marshal(bmarks)
	mockPluginAPI.On("KVGet", getBookmarksKey(userID)).Return(jsn, nil)
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
		want    int
		wantErr bool
	}{
		{
			name: "One Bookmark added",
			args: args{
				userID: "userID1",
			},
			wantErr: false,
			plugin: pluginWithMockedBookmarks("userID1", []*Bookmark{
				{
					PostID: "PostID_1",
					Title:  "Title_1",
				},
				{
					PostID: "PostID_2",
					Title:  "Title_2",
				},
			},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.plugin.addToBookmarksForUser(tt.args.userID, tt.args.newBookmark); (err != nil) != tt.wantErr {
				// bmarks, err := tt.plugin.getBookmarksForUser(tt.args.userID)
				// fmt.Printf("bmarks = %+v\n", bmarks)
				t.Errorf("Plugin.addToBookmarksForUser() error = %v, wantErr %v", err, tt.wantErr)
				// fmt.Printf("bmarks = %+v\n", bmarks)
			}
		})
	}
}
