package main

import (
	"encoding/json"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

// Bookmarks contains a map of bookmarks
type Bookmarks struct {
	ByID   map[string]*Bookmark
	api    plugin.API
	userID string
}

// NewBookmarksWithUser returns an initialized Labels for a User
func NewBookmarksWithUser(api plugin.API, userID string) *Bookmarks {
	return &Bookmarks{
		ByID:   make(map[string]*Bookmark),
		api:    api,
		userID: userID,
	}
}

func (b *Bookmarks) add(bmark *Bookmark) error {
	b.ByID[bmark.PostID] = bmark
	if err := b.storeBookmarks(); err != nil {
		return errors.Wrap(err, "failed to add bookmark")
	}
	return nil
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
	bmarkOrig.addLabelIDs(bmark.getLabelIDs())
	return bmark
}

// storeBookmarks stores all the users bookmarks
func (b *Bookmarks) storeBookmarks() error {
	jsonBookmarks, jsonErr := json.Marshal(b)
	if jsonErr != nil {
		return jsonErr
	}

	key := getBookmarksKey(b.userID)
	appErr := b.api.KVSet(key, jsonBookmarks)
	if appErr != nil {
		return appErr
	}

	return nil
}
