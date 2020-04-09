package main

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

const (
	// StoreBookmarksKey is the key used to store bookmarks in the plugin KV store
	StoreBookmarksKey = "bookmarks"
)

// storeBookmarks stores all the users bookmarks
func (p *Plugin) storeBookmarks(userID string, bmarks *Bookmarks) error {
	jsonBookmarks, jsonErr := json.Marshal(bmarks)
	if jsonErr != nil {
		return jsonErr
	}

	key := getBookmarksKey(userID)
	appErr := p.MattermostPlugin.API.KVSet(key, jsonBookmarks)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

// getBookmark returns a bookmark with the specified bookmarkID
func (p *Plugin) getBookmark(userID, bmarkID string) (*Bookmark, error) {
	bmarks, err := p.getBookmarks(userID)
	if err != nil {
		return nil, err
	}

	for _, bmark := range bmarks.ByID {
		if bmark.PostID == bmarkID {
			return bmark, nil
		}
	}

	return nil, nil
}

// addBookmark stores the bookmark in a map,
func (p *Plugin) addBookmark(userID string, bmark *Bookmark) (*Bookmarks, error) {

	// get all bookmarks for user
	bmarks, err := p.getBookmarks(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// no marks, initialize the store first
	if bmarks == nil {
		bmarks = NewBookmarks()
	}

	// user doesn't have any bookmarks add first bookmark and return
	if len(bmarks.ByID) == 0 {
		bmarks = NewBookmarks()
		bmarks.add(bmark)
		if err = p.storeBookmarks(userID, bmarks); err != nil {
			return nil, errors.New(err.Error())
		}
		return bmarks, nil
	}

	// bookmark already exists, update ModifiedAt and save
	if bmarks.exists(bmark.PostID) {
		bmarks.updateTimes(bmark.PostID)
		if err = p.storeBookmarks(userID, bmarks); err != nil {
			return nil, errors.New(err.Error())
		}
		return bmarks, nil
	}

	// bookmark doesn't exist. Add it
	bmarks.add(bmark)
	if err = p.storeBookmarks(userID, bmarks); err != nil {
		return nil, errors.New(err.Error())
	}
	return bmarks, nil
}

// getBookmarks returns unordered array of bookmarks for a user
func (p *Plugin) getBookmarks(userID string) (*Bookmarks, error) {
	originalJSONBookmarks, appErr := p.API.KVGet(getBookmarksKey(userID))
	if appErr != nil {
		return nil, appErr
	}

	if originalJSONBookmarks == nil {
		var bmarks *Bookmarks
		return bmarks, nil
	}

	var bmarks *Bookmarks
	jsonErr := json.Unmarshal(originalJSONBookmarks, &bmarks)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return bmarks, nil
}

// deleteBookmark deletes a bookmark from the store
func (p *Plugin) deleteBookmark(userID, bmarkID string) error {
	bmarks, err := p.getBookmarks(userID)
	if err != nil {
		return errors.New(err.Error())
	}

	if !bmarks.exists(bmarkID) {
		return errors.New(fmt.Sprintf("Bookmark `%v` does not exist", bmarkID))
	}

	bmarks.delete(bmarkID)
	p.storeBookmarks(userID, bmarks)

	return nil
}

func getBookmarksKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreBookmarksKey, userID)
}
