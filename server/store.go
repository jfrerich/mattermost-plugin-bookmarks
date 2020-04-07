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

// KVStore represents KVStore operations for bookmarks
type KVStore interface {
	storeBookmark(userID string, bmark *Bookmark) error
	storeBookmarks(userID string, bmarks *Bookmarks) error
	getBookmark(userID, bmarkID string) (*Bookmark, error)
	addBookmark(userID string, bmark *Bookmark) (*Bookmarks, error)
	deleteBookmark(userID, bmarkID string) error
	getBookmarks(userID string) (*Bookmarks, error)
	getBookmarksKey(userID string) string
}

type kvstore struct {
	plugin *Plugin
}

// NewStore creates a new kvstore
func NewStore(p *Plugin) *kvstore {
	return &kvstore{plugin: p}
}

// storeBookmark adds a bookmark for the user
func (s *kvstore) storeBookmark(userID string, bmark *Bookmark) error {
	_, err := s.addBookmark(userID, bmark)
	if err != nil {
		s.deleteBookmark(userID, bmark.PostID)
		return errors.New(err.Error())
	}

	return nil
}

// storeBookmarks stores all the users bookmarks
func (s *kvstore) storeBookmarks(userID string, bmarks *Bookmarks) error {
	jsonBookmarks, jsonErr := json.Marshal(bmarks)
	if jsonErr != nil {
		return jsonErr
	}

	key := getBookmarksKey(userID)
	appErr := s.plugin.API.KVSet(key, jsonBookmarks)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

// getBookmark returns a bookmark with the specified bookmarkID
func (s *kvstore) getBookmark(userID, bmarkID string) (*Bookmark, error) {
	bmarks, err := s.getBookmarks(userID)
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
func (s *kvstore) addBookmark(userID string, bmark *Bookmark) (*Bookmarks, error) {

	// get all bookmarks for user
	bmarks, err := s.getBookmarks(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// user doesn't have any bookmarks add first bookmark and return
	if len(bmarks.ByID) == 0 {
		bmarks = NewBookmarks()
		bmarks.add(bmark)
		if err = s.storeBookmarks(userID, bmarks); err != nil {
			return nil, errors.New(err.Error())
		}
		return bmarks, nil
	}

	// bookmark already exists, update ModifiedAt and save
	if bmarks.exists(bmark.PostID) {
		bmarks.updateTimes(bmark.PostID)
		if err = s.storeBookmarks(userID, bmarks); err != nil {
			return nil, errors.New(err.Error())
		}
		return bmarks, nil
	}

	// bookmark doesn't exist. Add it
	bmarks.add(bmark)
	if err = s.storeBookmarks(userID, bmarks); err != nil {
		return nil, errors.New(err.Error())
	}
	return bmarks, nil
}

// getBookmarks returns unordered array of bookmarks for a user
func (s *kvstore) getBookmarks(userID string) (*Bookmarks, error) {
	originalJSONBookmarks, appErr := s.plugin.API.KVGet(getBookmarksKey(userID))
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
func (s *kvstore) deleteBookmark(userID, bmarkID string) error {
	bmarks, err := s.getBookmarks(userID)
	if err != nil {
		return errors.New(err.Error())
	}

	if !bmarks.exists(bmarkID) {
		return errors.New(fmt.Sprintf("Bookmark `%v` does not exist", bmarkID))
	}

	bmarks.delete(bmarkID)
	s.storeBookmarks(userID, bmarks)

	return nil
}

func getBookmarksKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreBookmarksKey, userID)
}
