package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
)

// Bookmarks contains a map of bookmarks
type Bookmarks struct {
	ByID   map[string]*Bookmark
	Labels *Labels
}

// Bookmark contains information about an individual bookmark
type Bookmark struct {
	PostID     string   `json:"postid"`           // PostID is the ID for the bookmarked post and doubles as the Bookmark ID
	Title      string   `json:"title"`            // Title given to the bookmark
	CreateAt   int64    `json:"createAt"`         // The original creation time of the bookmark
	ModifiedAt int64    `json:"modifiedAt"`       // The original creation time of the bookmark
	LabelNames []string `json:"labels:omitempty"` // Array of labels added to the bookmark
}

// NewBookmarks returns an initialized Bookmarks struct
func NewBookmarks() *Bookmarks {
	bmarks := new(Bookmarks)
	bmarks.ByID = make(map[string]*Bookmark)
	bmarks.Labels = new(Labels)
	bmarks.Labels.ByName = make(map[string]*Label)
	return bmarks
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

func (b *Bookmark) hasUserTitle(bmark *Bookmark) bool {
	if bmark.Title != "" {
		return true
	}
	return false
}

func (b *Bookmark) hasLabels(bmark *Bookmark) bool {
	if bmark.LabelNames != nil {
		return true
	}
	return false
}

func (b *Bookmarks) labelExists(labelName string) (*Label, bool) {
	if label, ok := b.Labels.ByName[labelName]; ok {
		return label, true
	}
	return nil, false
}

func (b *Bookmarks) getLabel(labelName string) (*Label, bool) {
	if label, ok := b.Labels.ByName[labelName]; ok {
		return label, true
	}
	return nil, false
}

func (b *Bookmarks) addLabel(label *Label) {
	b.Labels.ByName[label.Name] = label
}

func (b *Bookmarks) deleteLabel(labelName string) {
	delete(b.Labels.ByName, labelName)
}
