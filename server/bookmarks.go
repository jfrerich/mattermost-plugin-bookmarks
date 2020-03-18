package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
)

// Bookmarks contains a map of bookmarks
type Bookmarks struct {
	ByID map[string]*Bookmark
}

// Bookmark contains information about an individual bookmark
type Bookmark struct {
	PostID     string  `json:"postid"`           // PostID is the ID for the bookmarked post and doubles as the Bookmark ID
	Title      string  `json:"title"`            // Title given to the bookmark
	CreateAt   int64   `json:"createAt"`         // The original creation time of the bookmark
	ModifiedAt int64   `json:"modifiedAt"`       // The original creation time of the bookmark
	Labels     []Label `json:"labels:omitempty"` // Array of labels added to the bookmark
}

// Label defines the parameters of a label
type Label struct {
	Name  string `json:"name"`
	Color string `json:"displayName"`
}

func (b *Bookmarks) add(bmark *Bookmark) {
	b.ByID[bmark.PostID] = bmark
}

func (b *Bookmarks) get(bmark *Bookmark) *Bookmark {
	return b.ByID[bmark.PostID]
}

func (b *Bookmarks) delete(bmarkID string) {
	delete(b.ByID, bmarkID)
}

func (b *Bookmarks) exists(bmark *Bookmark) bool {
	if _, ok := b.ByID[bmark.PostID]; ok {
		return true
	}
	return false
}

func (b *Bookmarks) updateTimes(bmark *Bookmark) *Bookmark {
	if bmark.CreateAt == 0 {
		bmark.CreateAt = model.GetMillis()
		bmark.ModifiedAt = bmark.CreateAt
	}
	bmark.ModifiedAt = model.GetMillis()
	return bmark
}

func (b *Bookmarks) new() *Bookmarks {
	bmarks := new(Bookmarks)
	bmarks.ByID = make(map[string]*Bookmark)
	return bmarks
}
