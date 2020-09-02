package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jfrerich/mattermost-plugin-bookmarks/server/bookmarks"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
)

const (
	UserID = "UserID"

	p1ID = "ID1"
	p2ID = "ID2"
	p3ID = "ID3"
	p4ID = "ID4"

	b1Title = "Title1 - New Bookmark - times are zero"
	b2Title = "Title2 - bookmarks initialized. Times created and same"
	b3Title = "Title3 - bookmarks already updated once"
)

func getHTTPTestBookmarks() *bookmarks.Bookmarks {
	b1 := &bookmarks.Bookmark{
		PostID:   p1ID,
		Title:    b1Title,
		LabelIDs: []string{"UUID1", "UUID2"},
	}
	b2 := &bookmarks.Bookmark{
		PostID:     p2ID,
		Title:      b2Title,
		CreateAt:   model.GetMillis() + 5,
		ModifiedAt: model.GetMillis(),
		LabelIDs:   []string{"UUID1", "UUID2"},
	}
	b3 := &bookmarks.Bookmark{
		PostID:     p3ID,
		Title:      b3Title,
		CreateAt:   model.GetMillis() + 2,
		ModifiedAt: model.GetMillis(),
	}
	b4 := &bookmarks.Bookmark{
		PostID:     p4ID,
		CreateAt:   model.GetMillis() + 3,
		ModifiedAt: model.GetMillis(),
	}

	bmarks := bookmarks.NewBookmarks(UserID)
	bmarks.ByID[b1.PostID] = b1
	bmarks.ByID[b2.PostID] = b2
	bmarks.ByID[b3.PostID] = b3
	bmarks.ByID[b4.PostID] = b4

	return bmarks
}

func getExecuteCommandTestLabels(t *testing.T) *bookmarks.Labels {
	l1 := &bookmarks.Label{
		Name: "label1",
	}
	l2 := &bookmarks.Label{
		Name: "label2",
	}
	l3 := &bookmarks.Label{
		Name: "label8",
	}

	labels := bookmarks.NewLabels(UserID)
	labels.ByID["UUID1"] = l1
	labels.ByID["UUID2"] = l2
	labels.ByID["UUID3"] = l3

	return labels
}

func TestHandleAddBookmark(t *testing.T) {
	b1 := bookmarks.Bookmark{
		Title:  "PostID-Title",
		PostID: "PostID1",
	}
	b2 := bookmarks.Bookmark{
		Title:  "PostID-Title",
		PostID: "ID1",
	}
	b3 := bookmarks.Bookmark{
		Title:    "PostID-Title",
		PostID:   "ID3",
		LabelIDs: []string{"newLabel"},
	}

	type bmarkWithChannel struct {
		Bookmark  *bookmarks.Bookmark `json:"bookmark"`
		ChannelID string              `json:"channelId"`
	}

	bmarks := getHTTPTestBookmarks()

	tests := map[string]struct {
		userID              string
		bookmark            *bookmarks.Bookmark
		bookmarks           *bookmarks.Bookmarks
		expectedCode        int
		expectedMsgPrefix   string
		expectedContains    []string
		expectedNotContains []string
	}{
		"Unauthed User": {
			bookmark:          &b1,
			bookmarks:         bmarks,
			expectedCode:      http.StatusUnauthorized,
			expectedMsgPrefix: "",
			expectedContains:  nil,
		},
		"Add first bookmark": {
			userID:            UserID,
			bookmark:          &b1,
			bookmarks:         bmarks,
			expectedCode:      http.StatusOK,
			expectedMsgPrefix: "Saved Bookmark",
			expectedContains: []string{
				fmt.Sprintf("[:link:](https://myhost.com/_redirect/pl/%v) **_PostID-Title_**", b1.PostID)},
		},
		"overwrite bookmark that exists": {
			userID:            UserID,
			bookmark:          &b2,
			bookmarks:         bmarks,
			expectedCode:      http.StatusOK,
			expectedMsgPrefix: "Saved Bookmark",
			expectedContains: []string{
				fmt.Sprintf("[:link:](https://myhost.com/_redirect/pl/%v) **_PostID-Title_**", b2.PostID)},
		},
		"bookmark contains labelID that is name of a bookmark": {
			userID:            UserID,
			bookmark:          &b3,
			bookmarks:         bmarks,
			expectedCode:      http.StatusOK,
			expectedMsgPrefix: "Saved Bookmark",
			expectedContains: []string{
				fmt.Sprintf("[:link:](https://myhost.com/_redirect/pl/%v) `newLabel` **_PostID-Title_**", b3.PostID)},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// ctrl := gomock.NewController(t)
			// defer ctrl.Finish()
			// mockPluginAPI := mock_pluginapi.NewMockAPI(ctrl)

			api := makeAPIMock()
			p := makePlugin(api)
			bWithChannel := bmarkWithChannel{
				Bookmark:  tt.bookmark,
				ChannelID: "SomeChannel",
			}
			jsonBmark, err := json.Marshal(bWithChannel)
			assert.Nil(t, err)
			jsonBmarks, err := json.Marshal(tt.bookmarks)
			assert.Nil(t, err)

			siteURL := "https://myhost.com"

			api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
			api.On("KVGet", bookmarks.GetBookmarksKey(UserID)).Return(jsonBmarks, nil)
			api.On("KVGet", bookmarks.GetLabelsKey(UserID)).Return(nil, nil)
			api.On("GetPost", tt.bookmark.PostID).Return(&model.Post{Message: "this is the post.Message"}, nil)

			api.On("GetConfig", mock.Anything).Return(&model.Config{ServiceSettings: model.ServiceSettings{SiteURL: &siteURL}})
			// mockPluginAPI.EXPECT().KVSet(bookmarks.GetBookmarksKey(UserID), gomock.Any)
			// mockPluginAPI.EXPECT().KVSet("bookmarks_UserID", jsonBmark).Return(nil)
			// mockPluginAPI.EXPECT().KVSet(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			// mockPluginAPI.EXPECT().KVGet(bookmarks.GetBookmarksKey(UserID)).Return(jsonBmarks, nil).AnyTimes()
			// mockPluginAPI.EXPECT().KVGet(bookmarks.GetLabelsKey(UserID)).Return(nil, nil).AnyTimes()
			// mockPluginAPI.EXPECT().GetPost(tt.bookmark.PostID).Return(&model.Post{Message: "this is the post.Message"}, nil).AnyTimes()
			// mockPluginAPI.EXPECT().GetConfig(mock.Anything).Return(&model.Config{ServiceSettings: model.ServiceSettings{SiteURL: &siteURL}})

			if tt.expectedCode == http.StatusOK {
				api.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.AnythingOfType("*model.Post")).Run(func(args mock.Arguments) {
					// 		// isSendEphemeralPostCalled = true
					//
					// 		post := args.Get(1).(*model.Post)
					// 		actual := strings.TrimSpace(post.Message)
					// 		assert.True(t, strings.HasPrefix(actual, tt.expectedMsgPrefix), "Expected returned message to start with: \n%s\nActual:\n%s", tt.expectedMsgPrefix, actual)
					// 		if tt.expectedContains != nil {
					// 			for i := range tt.expectedContains {
					// 				assert.Contains(t, actual, tt.expectedContains[i])
					// 			}
					// 		}
					// 		if tt.expectedNotContains != nil {
					// 			for i := range tt.expectedNotContains {
					// 				assert.NotContains(t, actual, tt.expectedNotContains[i])
					// 			}
					// 		}
					// 		// assert.Contains(t, actual, tt.expectedMsgPrefix)
				}).Once().Return(&model.Post{})
			}

			r := httptest.NewRequest(http.MethodPost, "/api/v1/add", strings.NewReader(string(jsonBmark)))
			r.Header.Add("Mattermost-User-Id", tt.userID)

			p.initialiseAPI()
			w := httptest.NewRecorder()
			p.ServeHTTP(nil, w, r)

			result := w.Result()
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedCode, result.StatusCode)
		})
	}
}

//nolint
func makePlugin(api *plugintest.API) *Plugin {
	p := &Plugin{}
	p.SetAPI(api)
	return p
}

func TestHandleGetBookmark(t *testing.T) {
	tests := map[string]struct {
		userID       string
		bookmark     *bookmarks.Bookmark
		bookmarks    *bookmarks.Bookmarks
		expectedCode int
	}{
		"Unauthed User": {
			expectedCode: http.StatusUnauthorized,
		},
		"get bookmark1": {
			userID:       UserID,
			expectedCode: http.StatusOK,
		},
	}
	for name, tt := range tests {
		api := makeAPIMock()
		p := makePlugin(api)

		bmarks := getHTTPTestBookmarks()
		bookmark := bmarks.ByID["ID1"]

		t.Run(name, func(t *testing.T) {
			jsonBmark, err := json.Marshal(bookmark)
			assert.Nil(t, err)
			jsonBmarks, err := json.Marshal(bmarks)
			assert.Nil(t, err)

			api.On("KVGet", bookmarks.GetBookmarksKey(UserID)).Return(jsonBmarks, nil)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/get?postID=ID1", strings.NewReader(string(jsonBmark)))
			r.Header.Add("Mattermost-User-Id", tt.userID)

			p.initialiseAPI()
			w := httptest.NewRecorder()
			p.ServeHTTP(&plugin.Context{}, w, r)

			result := w.Result()
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedCode, result.StatusCode)
		})
	}
}

func TestHandleViewBookmarks(t *testing.T) {
	bmarks := getHTTPTestBookmarks()

	tests := map[string]struct {
		userID       string
		bookmark     *bookmarks.Bookmark
		bookmarks    *bookmarks.Bookmarks
		expectedCode int
	}{
		"Unauthed User": {
			bookmark:     bmarks.ByID["ID1"],
			bookmarks:    bmarks,
			expectedCode: http.StatusUnauthorized,
		},
		"get bookmark1": {
			userID:       UserID,
			bookmark:     bmarks.ByID["ID1"],
			bookmarks:    bmarks,
			expectedCode: http.StatusOK,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			api := makeAPIMock()
			p := makePlugin(api)

			jsonBmarks, err := json.Marshal(tt.bookmarks)
			assert.Nil(t, err)

			siteURL := "https://myhost.com"

			api.On("KVGet", bookmarks.GetBookmarksKey(UserID)).Return(jsonBmarks, nil)
			api.On("KVGet", bookmarks.GetLabelsKey(UserID)).Return(nil, nil)

			api.On("GetConfig", mock.Anything).Return(&model.Config{ServiceSettings: model.ServiceSettings{SiteURL: &siteURL}})
			api.On("GetPost", mock.Anything).Return(&model.Post{Message: "this is the post.Message"}, nil)

			r := httptest.NewRequest(http.MethodPost, "/api/v1/view", strings.NewReader(string(jsonBmarks)))
			r.Header.Add("Mattermost-User-Id", tt.userID)

			if tt.expectedCode == http.StatusOK {
				api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(&model.Post{})
			}

			p.initialiseAPI()
			w := httptest.NewRecorder()
			p.ServeHTTP(&plugin.Context{}, w, r)

			result := w.Result()
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedCode, result.StatusCode)
		})
	}
}

func TestHandleLabelsGet(t *testing.T) {
	l1 := &bookmarks.Label{
		Name: "Label1",
		ID:   "LabelID1",
	}

	labels := getExecuteCommandTestLabels(t)
	tests := map[string]struct {
		userID       string
		label        *bookmarks.Label
		labels       *bookmarks.Labels
		expectedCode int
	}{
		"Unauthed User": {
			label:        l1,
			labels:       labels,
			expectedCode: http.StatusUnauthorized,
		},
		"No Errors": {
			userID:       UserID,
			label:        l1,
			labels:       labels,
			expectedCode: http.StatusOK,
		},
	}
	for name, tt := range tests {
		api := makeAPIMock()
		p := makePlugin(api)

		t.Run(name, func(t *testing.T) {
			jsonLabel, err := json.Marshal(tt.label)
			assert.Nil(t, err)

			api.On("KVGet", bookmarks.GetLabelsKey(UserID)).Return(jsonLabel, nil)

			r := httptest.NewRequest(http.MethodGet, "/api/v1/labels/get", strings.NewReader(string(jsonLabel)))
			r.Header.Add("Mattermost-User-Id", tt.userID)

			p.initialiseAPI()
			w := httptest.NewRecorder()
			p.ServeHTTP(nil, w, r)

			result := w.Result()
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedCode, result.StatusCode)
		})
	}
}

func TestHandleLabelsAdd(t *testing.T) {
	l1 := &bookmarks.Label{
		Name: "Label1",
		ID:   "LabelID1",
	}

	api := makeAPIMock()
	p := makePlugin(api)

	labels := getExecuteCommandTestLabels(t)
	tests := map[string]struct {
		userID       string
		label        *bookmarks.Label
		labels       *bookmarks.Labels
		expectedCode int
	}{
		"Unauthed User": {
			label:        l1,
			labels:       labels,
			expectedCode: http.StatusUnauthorized,
		},
		"Add a Label": {
			userID:       UserID,
			label:        l1,
			labels:       labels,
			expectedCode: http.StatusOK,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			jsonLabel, err := json.Marshal(tt.label)
			assert.Nil(t, err)

			siteURL := "https://myhost.com"
			api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
			api.On("KVGet", bookmarks.GetLabelsKey(UserID)).Return(nil, nil)
			api.On("GetConfig", mock.Anything).Return(&model.Config{ServiceSettings: model.ServiceSettings{SiteURL: &siteURL}})

			r := httptest.NewRequest(http.MethodPost, "/api/v1/labels/add?labelName=LabelID1", strings.NewReader(string(jsonLabel)))
			r.Header.Add("Mattermost-User-Id", tt.userID)

			p.initialiseAPI()
			w := httptest.NewRecorder()
			p.ServeHTTP(nil, w, r)

			result := w.Result()
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedCode, result.StatusCode)
		})
	}
}

func makeAPIMock() *plugintest.API {
	api := &plugintest.API{}

	api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
	api.On("LogWarn", mock.Anything, mock.Anything, mock.Anything).Maybe()
	api.On("LogError", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

	return api
}
