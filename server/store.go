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

// storeBookmarkForUser adds a bookmark for the user
func (p *Plugin) storeBookmarkForUser(userID string, bmark *Bookmark) error {
	_, err := p.addToBookmarksForUser(userID, bmark)
	fmt.Printf("err = %+v\n", err)
	if err != nil {
		p.deleteBookmarkForUser(userID, bmark.PostID)
		return errors.New(err.Error())
	}

	return nil
}

// storeBookmarks stores all the users bookmarks
func (p *Plugin) storeBookmarks(userID string, bmarks *Bookmarks) error {
	jsonBookmark, jsonErr := json.Marshal(bmarks)
	if jsonErr != nil {
		return jsonErr
	}

	appErr := p.API.KVSet(getBookmarksKey(userID), jsonBookmark)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

// getBookmark returns a bookmark with the specified bookmarkID
func (p *Plugin) getBookmark(userID, bmarkID string) (*Bookmark, error) {
	bmarks, err := p.getBookmarksForUser(userID)
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
func (p *Plugin) addToBookmarksForUser(userID string, bmark *Bookmark) (*Bookmarks, error) {
	bmarks, err := p.getBookmarksForUser(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// user doesn't have any bookmarks add first bookmark and return
	if bmarks == nil {
		bmarks := bmarks.new()
		bmarks.add(bmark)
		if err = p.storeBookmarks(userID, bmarks); err != nil {
			return nil, errors.New(err.Error())
		}
		return bmarks, nil
	}

	// bookmark already exists, update ModifiedAt and save
	if bmarks.exists(bmark) {
		// grab the saved bookmark from the store (includes original createdAt and
		// last modifiedAt times)
		bmark := bmarks.get(bmark)
		bmarks.updateTimes(bmark)
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

// getBookmarksForUser returns unordered array of bookmarks for a user
func (p *Plugin) getBookmarksForUser(userID string) (*Bookmarks, error) {
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

// deleteBookmarkForUser deletes a bookmark from the store
func (p *Plugin) deleteBookmarkForUser(userID, bmarkID string) error {
	bmarks, err := p.getBookmarksForUser(userID)
	if err != nil {
		return errors.New(err.Error())
	}

	// user doesn't have any bookmarks
	if bmarks == nil {
		return errors.New("User has no bookmarks")
	}

	bmarks.delete(bmarkID)
	p.storeBookmarks(userID, bmarks)

	return errors.New("unable to delete bookmark")
}

func getBookmarksKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreBookmarksKey, userID)
}
