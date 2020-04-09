package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleAdd(t *testing.T) {
	makePlugin := func(api *plugintest.API) *Plugin {
		p := &Plugin{}
		p.SetAPI(api)
		return p
	}
	api := makeAPIMock()
	p := makePlugin(api)

	t.Run("add bookmark", func(t *testing.T) {

		// get default bmark
		bmark := &Bookmark{
			Title:  "PostID-Title",
			PostID: "PostID1",
		}

		jsonBmark, err := json.Marshal(bmark)
		assert.Nil(t, err)

		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		api.On("KVGet", mock.Anything).Return(jsonBmark, nil)
		api.On("addBookmark", mock.Anything, mock.Anything).Return(jsonBmark, nil)

		r := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(string(jsonBmark)))
		r.Header.Add("Mattermost-User-Id", "theuserid")

		w := httptest.NewRecorder()
		p.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)
		assert.Equal(t, http.StatusOK, result.StatusCode)
	})
}
