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
	storeBookmarkForUser(userID string, bmark *Bookmark) error
	storeBookmarks(userID string, bmarks *Bookmarks) error
	getBookmark(userID, bmarkID string) (*Bookmark, error)
	addToBookmarksForUser(userID string, bmark *Bookmark) (*Bookmarks, error)
	deleteBookmarkForUser(userID, bmarkID string) error
	getBookmarksForUser(userID string) (*Bookmarks, error)
	getBookmarksKey(userID string) string
}

type kvstore struct {
	plugin *Plugin
}

// NewStore creates a new kvstore
func NewStore(p *Plugin) *kvstore {
	return &kvstore{plugin: p}
}

// storeBookmarkForUser adds a bookmark for the user
func (s *kvstore) storeBookmarkForUser(userID string, bmark *Bookmark) error {
	_, err := s.addToBookmarksForUser(userID, bmark)
	fmt.Printf("err = %+v\n", err)
	if err != nil {
		s.deleteBookmarkForUser(userID, bmark.PostID)
		return errors.New(err.Error())
	}

	return nil
}

// storeBookmarks stores all the users bookmarks
func (s *kvstore) storeBookmarks(userID string, bmarks *Bookmarks) error {
	jsonBookmark, jsonErr := json.Marshal(bmarks)
	if jsonErr != nil {
		return jsonErr
	}

	appErr := s.plugin.API.KVSet(getBookmarksKey(userID), jsonBookmark)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

// getBookmark returns a bookmark with the specified bookmarkID
func (s *kvstore) getBookmark(userID, bmarkID string) (*Bookmark, error) {
	bmarks, err := s.getBookmarksForUser(userID)
	if err != nil {
		fmt.Printf("Error = %+v\n", err)
	}

	for _, bmark := range bmarks.ByID {
		if bmark.PostID == bmarkID {
			return bmark, nil
		}
	}

	return nil, nil
}

// addToBookmarksForUser stores the bookmark in a map,
func (s *kvstore) addToBookmarksForUser(userID string, bmark *Bookmark) (*Bookmarks, error) {
	bmarks, err := s.getBookmarksForUser(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// user doesn't have any bookmarks add first bookmark and return
	if bmarks == nil {
		bmarks = NewBookmarks()
		bmarks.add(bmark)
		if err = s.storeBookmarks(userID, bmarks); err != nil {
			return nil, errors.New(err.Error())
		}
		return bmarks, nil
	}

	// bookmark already exists, update ModifiedAt and save
	if bmarks.exists(bmark.PostID) {
		// grab the saved bookmark from the store (includes original createdAt and
		// last modifiedAt times)
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

// getBookmarksForUser returns unordered array of bookmarks for a user
func (s *kvstore) getBookmarksForUser(userID string) (*Bookmarks, error) {
	key := getBookmarksKey(userID)

	fmt.Printf("userID = %+v\n", userID)
	fmt.Printf("key = %+v\n", key)

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

// deleteBookmarkForUser deletes a bookmark from the store
func (s *kvstore) deleteBookmarkForUser(userID, bmarkID string) error {
	bmarks, err := s.getBookmarksForUser(userID)
	if err != nil {
		return errors.New(err.Error())
	}

	// user doesn't have any bookmarks
	if bmarks == nil {
		return errors.New("User has no bookmarks")
	}

	bmarks.delete(bmarkID)
	s.storeBookmarks(userID, bmarks)

	return errors.New("unable to delete bookmark")
}

func getBookmarksKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreBookmarksKey, userID)
}
