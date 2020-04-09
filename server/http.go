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
	case "/view":
		p.handleView(w, r)
	case "/delete":
		p.handleDelete(w, r)
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

	_, err = p.addBookmark(userID, bmark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (p *Plugin) handleDelete(w http.ResponseWriter, r *http.Request) {
	return
}

func (p *Plugin) handleView(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var bmark Bookmark
	fmt.Printf("r.Body = %+v\n", r.Body)
	decoder := json.NewDecoder(r.Body)
	decoderD, _ := json.MarshalIndent(decoder, "", "    ")
	fmt.Printf("decoder = %+v\n", string(decoderD))
	err := decoder.Decode(&bmark)
	if err != nil {
		p.API.LogError("Unable to decode JSON err=" + err.Error())
		p.handleErrorWithCode(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	p.addBookmark(userID, &bmark)

	return
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
