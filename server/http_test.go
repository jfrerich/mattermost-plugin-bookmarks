package main

import (
	"testing"
)

func TestHandleAdd(t *testing.T) {
	// makePlugin := func(api *plugintest.API) *Plugin {
	// 	p := &Plugin{}
	// 	p.SetAPI(api)
	// 	return p
	// }

	// t.Run("add bookmark", func(t *testing.T) {
	// 	api := makeAPIMock()
	// 	p := makePlugin(api)
	//
	// 	// get default bmark
	// 	bmark := &Bookmark{
	// 		Title:  "PostID-Title",
	// 		PostID: "PostID1",
	// 	}
	//
	// 	jsonBmark, err := json.Marshal(bmark)
	// 	assert.Nil(t, err)
	//
	// 	api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
	// 	api.On("KVGet", getBookmarksKey("theuserid")).Return(nil, nil)
	//
	// 	r := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(string(jsonBmark)))
	// 	r.Header.Add("Mattermost-User-Id", "theuserid")
	//
	// 	w := httptest.NewRecorder()
	// 	p.ServeHTTP(nil, w, r)
	//
	// 	result := w.Result()
	// 	assert.NotNil(t, result)
	// 	assert.Equal(t, http.StatusOK, result.StatusCode)
	//
	// 	api.On("KVGet", getBookmarksKey("theuserid2")).Return(jsonBmark, nil)
	//
	// 	r2 := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(string(jsonBmark)))
	// 	r2.Header.Add("Mattermost-User-Id", "theuserid2")
	//
	// 	w2 := httptest.NewRecorder()
	// 	p.ServeHTTP(nil, w2, r2)
	//
	// 	result2 := w2.Result()
	// 	assert.NotNil(t, result2)
	// 	assert.Equal(t, http.StatusOK, result2.StatusCode)
	//
	// })
}
