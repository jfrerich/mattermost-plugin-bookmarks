package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	api := makeAPIMock()
	p := makePlugin(api)

	bmarks := getExecuteCommandTestBookmarks()

	tests := map[string]struct {
		userID       string
		bookmark     *Bookmark
		bookmarks    *Bookmarks
		expectedCode int
	}{
		"Unauthed User": {
			bookmark:     b1,
			bookmarks:    bmarks,
			expectedCode: http.StatusUnauthorized,
		},
		"Add first bookmark": {
			userID:       UserID,
			bookmark:     b1,
			bookmarks:    bmarks,
			expectedCode: http.StatusOK,
		},
		"overwrite bookmark that exists": {
			userID:       UserID,
			bookmark:     b2,
			bookmarks:    bmarks,
			expectedCode: http.StatusOK,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			jsonBmark, err := json.Marshal(tt.bookmark)
			assert.Nil(t, err)
			jsonBmarks, err := json.Marshal(tt.bookmarks)
			assert.Nil(t, err)

			api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
			api.On("KVGet", getBookmarksKey(UserID)).Return(jsonBmarks, nil)
			api.On("KVGet", getLabelsKey(UserID)).Return(nil, nil)

			r := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(string(jsonBmark)))
			r.Header.Add("Mattermost-User-Id", tt.userID)

			w := httptest.NewRecorder()
			p.ServeHTTP(nil, w, r)

			result := w.Result()
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedCode, result.StatusCode)
		})
	}
}
