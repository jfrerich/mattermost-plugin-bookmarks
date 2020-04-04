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

	// var bmarks *Bookmarks
	bmarks := Bookmarks{ByID: map[string]*Bookmark{}}
	bmarks.ByID[userID] = bookmark
	// bmarks.ByID["lj"] = bookmarks
	// bmarksD, _ := json.MarshalIndent(bmarks, "", "    ")
	// fmt.Printf("bmarks = %+v\n", string(bmarksD))

	// jsn, _ := json.Marshal(bmarks)
	mockPluginAPI.On("KVGet", getBookmarksKey(userID)).Return(nil, nil)
	mockPluginAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)

	// mockPluginAPI.On("KVSet", getBookmarksKey(userID)).Return(nil, nil)
	p.SetAPI(mockPluginAPI)

	return p
}

// wantedSubscriptions returns what should be returned after sorting by repo names
func wantedBookmarks(userID1 string, bookmarks []*Bookmark) *Bookmarks {
	var bmarks *Bookmarks
	bmarks.new()

	for k, v := range bookmarks {
		fmt.Printf("k = %+v\n", k)
		fmt.Printf("v = %+v\n", v)
		// bmarks.add(v)
		bmarks.ByID[userID1].PostID = v.PostID
		bmarks.ByID[userID1].Title = v.Title
		// subs = append(subs, &Subscription{
		// 	ChannelID:  chanelID,
		// 	Repository: st,
		// })
	}
	return bmarks
}

func TestPlugin_getBookmarksForUser(t *testing.T) {
	plugin := Plugin{}
	api := &plugintest.API{}
	plugin.SetAPI(api)

	api.On("KVGet", "bookmarks_randomUser").Return(nil, nil)
	bmarks, _ := plugin.kvstore.getBookmarksForUser("randomUser")
	fmt.Printf("bmarks = %+v\n", bmarks)

}

// plugin.SetHelpers(helpers)
// 	var b1 *Bookmark
// 	b1 = &Bookmark{PostID: "ID1", Title: "Title1"}
// helpers := &plugintest.Helpers{}
// helpers.On("KVGet", "randomUser").Return(nil, nil)

// type pluginTest struct {
// 	api    plugin.API
// 	plugin *Plugin
// }
//
// p := &pluginTest{}
// _, _ = p.plugin.getBookmarksForUser("randomUser")
//

// assert.Nil(err)

// bmarks, _ := ptest.getBookmarksForUser("userID1")
// fmt.Printf("bmarks = %+v\n", bmarks)
// assert.NotZero(t, int64(bmarks.ByID[getBookmarksKey("userID1")].CreateAt))

func TestStoreBookmarks(t *testing.T) {

	// intialize test Bookmarks
	u1 := "userID1"
	// u2 := "userID2"

	b1 := &Bookmark{PostID: "ID1", Title: "Title1"}
	b2 := &Bookmark{PostID: "ID2", Title: "Title2"}
	// b3 := &Bookmark{PostID: "ID3", Title: "Title3"}

	api := &plugintest.API{}

	// Add Bookmarks
	bmarks := Bookmarks{}
	bmarks = *bmarks.new()
	bmarks.add(b1)
	bmarks.add(b2)

	bmarksD, _ := json.MarshalIndent(bmarks, "", "    ")
	fmt.Printf("bmarks = %+v\n", string(bmarksD))

	var plugin Plugin
	plugin.SetAPI(api)

	// Markshal the bmarks and mock api call
	jsonBookmarks, err := json.Marshal(bmarks)
	api.On("KVSet", "bookmarks_userID1", jsonBookmarks).Return(nil)

	// store bmarks using API
	err = plugin.kvstore.storeBookmarks(u1, &bmarks)
	assert.Nil(t, err)

	jsonBookmarksD, _ := json.MarshalIndent(jsonBookmarks, "", "    ")
	fmt.Printf("jsonBookmarks = %+v\n", string(jsonBookmarksD))

	api.On("KVGet", "bookmarks_userID1").Return(jsonBookmarks, nil)
	getBmarks, err := plugin.kvstore.getBookmarksForUser(u1)

	fmt.Printf("getBmarks = %+v\n", getBmarks)

	assert.Equal(t, &bmarks, getBmarks)

}

// func TestPlugin_addToBookmarksForUser2(t *testing.T) {
// 	plugin := Plugin{}
// 	api := &plugintest.API{}
// 	plugin.SetAPI(api)
// 	var b1 *Bookmark
// 	b1 = &Bookmark{PostID: "ID1", Title: "Title1"}
// 	fmt.Printf("b1 = %+v\n", b1)
// 	api.On("KVGet", "bookmarks_userID1").Return(nil, nil)
// 	api.On("KVSet", "bookmarks_userID1").Return(nil)
// 	bmarks, _ := plugin.addToBookmarksForUser("userID1", b1)
// 	fmt.Printf("bmarks = %+v\n", bmarks)
// 	// 	// assert.NotZero(t, int64(bmarks.ByID[getBookmarksKey("userID1")].CreateAt))
// 	//
// }

// func TestPlugin_addToBookmarksForUser(t *testing.T) {
//
// 	// var bookmarks  *Bookmarks
// 	// bookmarks.ByID["jason"] =
//
// 	type args struct {
// 		userID      string
// 		newBookmark *Bookmark
// 	}
// 	tests := []struct {
// 		name    string
// 		plugin  *Plugin
// 		args    args
// 		want    *Bookmark
// 		wantErr bool
// 	}{
// 		{
// 			name: "One Bookmark added",
// 			args: args{
// 				userID: "userID1",
// 				newBookmark: &Bookmark{
// 					PostID: "PostID_1",
// 					Title:  "Title_1",
// 				},
// 			},
// 			wantErr: true,
// 			// plugin:  pluginWithMockedBookmarks("userID1"), &Bookmarks{
// 			// 	PostID: "PostID_1",
// 			// 	Title:  "Title_1",
// 			// },
//
// 			plugin: pluginWithMockedBookmarks("userID1",
// 				&Bookmark{
// 					PostID: "PostID_1",
// 					Title:  "Title_1",
// 				},
// 			),
// 			want: wantedBookmarks("userID1", []*Bookmark{"POSTID1": {PostID: "POSTID1", Title: "Title_1"}}),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got, err := tt.plugin.addToBookmarksForUser(tt.args.userID, tt.args.newBookmark); (err != nil) != tt.wantErr {
// 				gotD, _ := json.MarshalIndent(got, "", "    ")
// 				fmt.Printf("got = %+v\n", string(gotD))
// 				// bmarks, err := tt.plugin.getBookmarksForUser(tt.args.userID)
// 				// fmt.Printf("bmarks = %+v\n", bmarks)
// 				t.Errorf("Plugin.addToBookmarksForUser() error = %v, wantErr %v", err, tt.wantErr)
// 				assert.Equal(t, tt.want, got, "they should be the same")
// 				// fmt.Printf("bmarks = %+v\n", bmarks)
// 			}
// 		})
// 	}
// }
