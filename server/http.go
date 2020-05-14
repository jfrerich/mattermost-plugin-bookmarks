package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
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
	case "/labels/add":
		p.handleLabelsAdd(w, r)
	// case "/delete":
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) handleAdd(w http.ResponseWriter, r *http.Request) {
	type bmarkWithChannel struct {
		Bookmark  *Bookmark `json:"bookmark"`
		ChannelID string    `json:"channelId"`
	}

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

	var req *bmarkWithChannel
	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bmark := req.Bookmark
	channelID := req.ChannelID

	bmarks, err := NewBookmarksWithUser(p.API, userID).getBookmarks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	l, err := NewLabelsWithUser(p.API, userID).getLabels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	ids := bmark.getLabelIDs()

	var newIDs []string
	var label *Label
	for _, id := range ids {
		label, err = l.get(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// if doesn't exist, this is a name and needs to be added to the labels
		// store.  also save the id to the bookmark, not the name
		if label == nil {
			var labelNew *Label
			labelNew, err = l.addLabel(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			newIDs = append(newIDs, labelNew.ID)
			continue
		}
		newIDs = append(newIDs, id)
	}

	// update bmark with UUID values, not the names
	bmark.LabelIDs = newIDs
	err = bmarks.addBookmark(bmark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var names []string
	var name string
	for _, id := range newIDs {
		name, err = l.getNameFromID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		names = append(names, name)
	}

	text, err := p.getBmarkTextOneLine(bmark, names, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	message := "Saved Bookmark:\n" + text

	post := &model.Post{
		UserId:    p.getBotID(),
		ChannelId: channelID,
		Message:   message,
	}
	_ = p.API.SendEphemeralPost(userID, post)
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

	bmarks, err := NewBookmarksWithUser(p.API, userID).getBookmarks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// return nil if bookmark does not exist
	bmark, err := bmarks.getBookmark(postID)
	if bmark == nil {
		var bb []byte
		_, err = w.Write(bb)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	if err != nil {
		p.handleErrorWithCode(w, http.StatusBadRequest, "Unable to get bookmark", err)
	}

	resp, err := json.Marshal(bmark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *Plugin) handleLabelsAdd(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query()
	labelName := query["labelName"][0]
	l := NewLabelsWithUser(p.API, userID)
	labels, err := l.getLabels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	label, err := labels.addLabel(labelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(label)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
