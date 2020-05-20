package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleAdd(t *testing.T) {
	b1 := &Bookmark{
		Title:  "PostID-Title",
		PostID: "PostID1",
	}
	b2 := &Bookmark{
		Title:  "PostID-Title",
		PostID: "ID1",
	}
	b3 := &Bookmark{
		Title:    "PostID-Title",
		PostID:   "ID3",
		LabelIDs: []string{"newLabel"},
	}

	type bmarkWithChannel struct {
		Bookmark  *Bookmark `json:"bookmark"`
		ChannelID string    `json:"channelId"`
	}

	api := makeAPIMock()
	p := makePlugin(api)

	bmarks := getExecuteCommandTestBookmarks()

	tests := map[string]struct {
		userID              string
		bookmark            *Bookmark
		bookmarks           *Bookmarks
		expectedCode        int
		expectedMsgPrefix   string
		expectedContains    []string
		expectedNotContains []string
	}{
		"Unauthed User": {
			bookmark:          b1,
			bookmarks:         bmarks,
			expectedCode:      http.StatusUnauthorized,
			expectedMsgPrefix: "",
			expectedContains:  nil,
		},
		"Add first bookmark": {
			userID:            UserID,
			bookmark:          b1,
			bookmarks:         bmarks,
			expectedCode:      http.StatusOK,
			expectedMsgPrefix: "Saved Bookmark",
			expectedContains: []string{
				fmt.Sprintf("[:link:](https://myhost.com/_redirect/pl/%v) PostID-Title", b1.PostID)},
		},
		"overwrite bookmark that exists": {
			userID:            UserID,
			bookmark:          b2,
			bookmarks:         bmarks,
			expectedCode:      http.StatusOK,
			expectedMsgPrefix: "Saved Bookmark",
			expectedContains: []string{
				fmt.Sprintf("[:link:](https://myhost.com/_redirect/pl/%v) PostID-Title", b2.PostID)},
		},
		"bookmark contains labelID that is name of a bookmark": {
			userID:            UserID,
			bookmark:          b3,
			bookmarks:         bmarks,
			expectedCode:      http.StatusOK,
			expectedMsgPrefix: "Saved Bookmark",
			expectedContains: []string{
				fmt.Sprintf("[:link:](https://myhost.com/_redirect/pl/%v) `newLabel` PostID-Title", b3.PostID)},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
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
			api.On("KVGet", getBookmarksKey(UserID)).Return(jsonBmarks, nil)
			api.On("KVGet", getLabelsKey(UserID)).Return(nil, nil)
			api.On("GetPost", tt.bookmark.PostID).Return(&model.Post{Message: "this is the post.Message"}, nil)
			api.On("GetConfig", mock.Anything).Return(&model.Config{ServiceSettings: model.ServiceSettings{SiteURL: &siteURL}})

			if tt.expectedCode == http.StatusOK {
				api.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.AnythingOfType("*model.Post")).Run(func(args mock.Arguments) {
					// isSendEphemeralPostCalled = true

					post := args.Get(1).(*model.Post)
					actual := strings.TrimSpace(post.Message)
					assert.True(t, strings.HasPrefix(actual, tt.expectedMsgPrefix), "Expected returned message to start with: \n%s\nActual:\n%s", tt.expectedMsgPrefix, actual)
					if tt.expectedContains != nil {
						for i := range tt.expectedContains {
							assert.Contains(t, actual, tt.expectedContains[i])
						}
					}
					if tt.expectedNotContains != nil {
						for i := range tt.expectedNotContains {
							assert.NotContains(t, actual, tt.expectedNotContains[i])
						}
					}
					// assert.Contains(t, actual, tt.expectedMsgPrefix)
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

func TestHandleLabelsGet(t *testing.T) {
	l1 := &Label{
		Name: "Label1",
		ID:   "LabelID1",
	}

	api := makeAPIMock()
	p := makePlugin(api)

	labels := getExecuteCommandTestLabels()
	tests := map[string]struct {
		userID              string
		label               *Label
		labels              *Labels
		expectedCode        int
		expectedMsgPrefix   string
		expectedContains    []string
		expectedNotContains []string
	}{
		"Unauthed User": {
			label:             l1,
			labels:            labels,
			expectedCode:      http.StatusUnauthorized,
			expectedMsgPrefix: "",
			expectedContains:  nil,
		},
		"No Errors": {
			userID:            UserID,
			label:             l1,
			labels:            labels,
			expectedCode:      http.StatusOK,
			expectedMsgPrefix: "",
			expectedContains:  nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			jsonLabel, err := json.Marshal(tt.label)
			assert.Nil(t, err)

			siteURL := "https://myhost.com"
			api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
			api.On("KVGet", getLabelsKey(UserID)).Return(nil, nil)
			// api.On("GetPost", tt.bookmark.PostID).Return(&model.Post{Message: "this is the post.Message"}, nil)
			api.On("GetConfig", mock.Anything).Return(&model.Config{ServiceSettings: model.ServiceSettings{SiteURL: &siteURL}})

			if tt.expectedCode == http.StatusOK {
				api.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.AnythingOfType("*model.Post")).Run(func(args mock.Arguments) {
					// isSendEphemeralPostCalled = true

					post := args.Get(1).(*model.Post)
					actual := strings.TrimSpace(post.Message)
					assert.True(t, strings.HasPrefix(actual, tt.expectedMsgPrefix), "Expected returned message to start with: \n%s\nActual:\n%s", tt.expectedMsgPrefix, actual)
					if tt.expectedContains != nil {
						for i := range tt.expectedContains {
							assert.Contains(t, actual, tt.expectedContains[i])
						}
					}
					if tt.expectedNotContains != nil {
						for i := range tt.expectedNotContains {
							assert.NotContains(t, actual, tt.expectedNotContains[i])
						}
					}
					// assert.Contains(t, actual, tt.expectedMsgPrefix)
				}).Once().Return(&model.Post{})
			}

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
	l1 := &Label{
		Name: "Label1",
		ID:   "LabelID1",
	}

	api := makeAPIMock()
	p := makePlugin(api)

	labels := getExecuteCommandTestLabels()
	tests := map[string]struct {
		userID              string
		label               *Label
		labels              *Labels
		expectedCode        int
		expectedMsgPrefix   string
		expectedContains    []string
		expectedNotContains []string
	}{
		"Unauthed User": {
			label:             l1,
			labels:            labels,
			expectedCode:      http.StatusUnauthorized,
			expectedMsgPrefix: "",
			expectedContains:  nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			jsonLabel, err := json.Marshal(tt.label)
			assert.Nil(t, err)

			siteURL := "https://myhost.com"
			api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
			api.On("KVGet", getLabelsKey(UserID)).Return(nil, nil)
			// api.On("GetPost", tt.bookmark.PostID).Return(&model.Post{Message: "this is the post.Message"}, nil)
			api.On("GetConfig", mock.Anything).Return(&model.Config{ServiceSettings: model.ServiceSettings{SiteURL: &siteURL}})

			if tt.expectedCode == http.StatusOK {
				api.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.AnythingOfType("*model.Post")).Run(func(args mock.Arguments) {
					// isSendEphemeralPostCalled = true

					post := args.Get(1).(*model.Post)
					actual := strings.TrimSpace(post.Message)
					assert.True(t, strings.HasPrefix(actual, tt.expectedMsgPrefix), "Expected returned message to start with: \n%s\nActual:\n%s", tt.expectedMsgPrefix, actual)
					if tt.expectedContains != nil {
						for i := range tt.expectedContains {
							assert.Contains(t, actual, tt.expectedContains[i])
						}
					}
					if tt.expectedNotContains != nil {
						for i := range tt.expectedNotContains {
							assert.NotContains(t, actual, tt.expectedNotContains[i])
						}
					}
					// assert.Contains(t, actual, tt.expectedMsgPrefix)
				}).Once().Return(&model.Post{})
			}

			r := httptest.NewRequest(http.MethodPost, "/api/v1/labels/add", strings.NewReader(string(jsonLabel)))
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
