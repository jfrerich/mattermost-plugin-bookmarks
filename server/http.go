package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	case "/add":
		p.handleAdd(w, r)
	case "/get":
		p.handleView(w, r)
	// case "/delete":
	// 	p.handleDelete(w, r)
	case "/labels/get":
		p.handleLabelsGet(w, r)
	// case "/delete":
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) handleAdd(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var bmark *Bookmark
	if err = json.Unmarshal(body, &bmark); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b := NewBookmarksWithUser(p.API, userID)
	bmarks, err := b.getBookmarks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// ok if bmarks are nil.  They dont yet exist
	// if bmarks == nil {
	// 	// http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	fmt.Println("IN HERE")
	// 	http.Error(w, "no bookmarks", http.StatusInternalServerError)
	// 	// p.handleErrorWithCode(w, http.StatusBadRequest, "Unable to decode JSON", err)
	// 	return
	// }
	err = bmarks.addBookmark(bmark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// func (p *Plugin) handleDelete(w http.ResponseWriter, r *http.Request) {
// 	return
// }

func (p *Plugin) handleView(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query()
	postID := query["postID"][0]
	fmt.Printf("postID = %+v\n", postID)

	b := NewBookmarksWithUser(p.API, userID)
	bmarks, err := b.getBookmarks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	bmark, err := bmarks.getBookmark(postID)
	fmt.Printf("bmark = %+v\n", bmark)
	if err != nil {
		p.handleErrorWithCode(w, http.StatusBadRequest, "Unable to get bookmark", err)
	}

	resp, err := json.Marshal(bmark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)

}

func (p *Plugin) handleLabelsGet(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	l := NewLabelsWithUser(p.API, userID)
	labels, err := l.getLabels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resp, err := json.Marshal(labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (p *Plugin) handleErrorWithCode(w http.ResponseWriter, code int, errTitle string, err error) {
	w.WriteHeader(code)
	b, _ := json.Marshal(struct {
		Error   string `json:"error"`
		Details string `json:"details"`
	}{
		Error:   errTitle,
		Details: err.Error(),
	})
	_, _ = w.Write(b)
}
