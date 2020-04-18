package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Bookmarks contains a map of bookmarks
type Bookmarks struct {
	ByID map[string]*Bookmark
	api  plugin.API
}

// Bookmark contains information about an individual bookmark
type Bookmark struct {
	PostID     string   `json:"postid"`           // PostID is the ID for the bookmarked post and doubles as the Bookmark ID
	Title      string   `json:"title"`            // Title given to the bookmark
	CreateAt   int64    `json:"createAt"`         // The original creation time of the bookmark
	ModifiedAt int64    `json:"modifiedAt"`       // The original creation time of the bookmark
	LabelIDs   []string `json:"labels:omitempty"` // Array of labels added to the bookmark
}

// NewBookmarks returns an initialized Bookmarks struct
func NewBookmarks(api plugin.API) *Bookmarks {
	return &Bookmarks{
		ByID: make(map[string]*Bookmark),
		api:  api,
	}
}

func (b *Bookmarks) add(bmark *Bookmark) {
	b.ByID[bmark.PostID] = bmark
}

func (b *Bookmarks) get(bmarkID string) *Bookmark {
	return b.ByID[bmarkID]
}

func (b *Bookmarks) delete(bmarkID string) {
	delete(b.ByID, bmarkID)
}

func (b *Bookmarks) exists(bmarkID string) (*Bookmark, bool) {
	if bmark, ok := b.ByID[bmarkID]; ok {
		return bmark, true
	}
	return nil, false
}

func (b *Bookmarks) updateTimes(bmarkID string) *Bookmark {
	bmark := b.get(bmarkID)
	if bmark.CreateAt == 0 {
		bmark.CreateAt = model.GetMillis()
		bmark.ModifiedAt = bmark.CreateAt
	}
	bmark.ModifiedAt = model.GetMillis()
	return bmark
}

func (b *Bookmarks) updateLabels(bmark *Bookmark) *Bookmark {
	bmarkOrig := b.get(bmark.PostID)
	bmarkOrig.setLabelIDs(bmark.getLabelIDs())
	return bmark
}

func (b *Bookmark) hasUserTitle(bmark *Bookmark) bool {
	if bmark.getTitle() != "" {
		return true
	}
	return false
}

func (b *Bookmark) hasLabels(bmark *Bookmark) bool {
	if bmark.getLabelIDs() != nil {
		return true
	}
	return false
}

func (b *Bookmark) getTitle() string {
	return b.Title
}

func (b *Bookmark) setTitle(title string) {
	b.Title = title
	return
}

func (b *Bookmark) getLabelIDs() []string {
	return b.LabelIDs
}

func (b *Bookmark) setLabelIDs(IDs []string) {
	b.LabelIDs = IDs
	return
}
