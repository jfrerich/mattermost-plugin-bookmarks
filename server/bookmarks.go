package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

const (
	// StoreBookmarksKey is the key used to store bookmarks in the plugin KV store
	StoreBookmarksKey = "bookmarks"
)

// storeBookmarks stores all the users bookmarks
func (b *Bookmarks) storeBookmarks(userID string) error {
	jsonBookmarks, jsonErr := json.Marshal(b)
	if jsonErr != nil {
		return jsonErr
	}

	key := getBookmarksKey(userID)
	appErr := b.api.KVSet(key, jsonBookmarks)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

// getBookmark returns a bookmark with the specified bookmarkID
func (b *Bookmarks) getBookmark(userID, bmarkID string) (*Bookmark, error) {

	_, ok := b.exists(bmarkID)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Bookmark `%v` does not exist", bmarkID))
	}

	for _, bmark := range b.ByID {
		if bmark.PostID == bmarkID {
			return bmark, nil
		}
	}

	return nil, nil
}

// addBookmark stores the bookmark in a map,
func (b *Bookmarks) addBookmark(userID string, bmark *Bookmark) error {

	// user doesn't have any bookmarks add first bookmark and return
	if len(b.ByID) == 0 {
		b.add(bmark)
		if err := b.storeBookmarks(userID); err != nil {
			return errors.New(err.Error())
		}
		return nil
	}

	// bookmark already exists, update ModifiedAt and save
	_, ok := b.exists(bmark.PostID)
	if ok {
		b.updateTimes(bmark.PostID)
		b.updateLabels(bmark)

		if err := b.storeBookmarks(userID); err != nil {
			return errors.New(err.Error())
		}
		return nil
	}

	// bookmark doesn't exist. Add it
	b.add(bmark)
	if err := b.storeBookmarks(userID); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// getBookmarks returns a users bookmarks.  If the user has no bookmarks,
// return nil bookmarks
func (b *Bookmarks) getBookmarks(userID string) (*Bookmarks, error) {

	// if a user not not have bookmarks, bb will be nil
	bb, appErr := b.api.KVGet(getBookmarksKey(userID))
	if appErr != nil {
		return nil, appErr
	}

	if bb == nil {
		return nil, nil
	}

	// return initialized bookmarks
	bmarks := NewBookmarks(b.api)
	jsonErr := json.Unmarshal(bb, &bmarks)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return bmarks, nil
}

// ByPostCreateAt returns an array of bookmarks sorted by post.CreateAt times
func (b *Bookmarks) ByPostCreateAt(bmarks *Bookmarks) ([]*Bookmark, error) {
	// build temp map
	tempMap := make(map[int64]string)
	for _, bmark := range bmarks.ByID {
		post, appErr := b.api.GetPost(bmark.PostID)
		if appErr != nil {
			return nil, appErr
		}
		tempMap[post.CreateAt] = bmark.PostID
	}

	// sort post.CreateAt (keys)
	keys := make([]int, 0, len(tempMap))
	for k := range tempMap {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	// reconstruct the bookmarks in a sorted array
	var bookmarks []*Bookmark
	for _, k := range keys {
		bmark := bmarks.ByID[tempMap[int64(k)]]
		bookmarks = append(bookmarks, bmark)
	}

	return bookmarks, nil
}

func (b *Bookmarks) getBookmarksWithLabelID(userID, labelID string) (*Bookmarks, error) {

	bmarksWithLabel := NewBookmarks(b.api)

	for _, bmark := range b.ByID {
		if bmark.hasLabels(bmark) {
			for _, id := range bmark.getLabelIDs() {
				if labelID == id {
					bmarksWithLabel.add(bmark)
				}
			}
		}
	}

	return bmarksWithLabel, nil
}

// deleteBookmark deletes a bookmark from the store
func (b *Bookmarks) deleteBookmark(userID, bmarkID string) (*Bookmark, error) {
	var bmark *Bookmark

	_, ok := b.exists(bmarkID)
	if !ok {
		return bmark, errors.New(fmt.Sprintf("Bookmark `%v` does not exist", bmarkID))
	}

	bmark = b.get(bmarkID)

	b.delete(bmarkID)
	b.storeBookmarks(userID)

	return bmark, nil
}

// deleteLabel deletes a label from a bookmark
func (b *Bookmarks) deleteLabel(userID, bmarkID string, labelID string) error {
	bmark, err := b.getBookmark(userID, bmarkID)
	if err != nil {
		return errors.New(err.Error())
	}

	origLabels := bmark.getLabelIDs()

	var newLabels []string
	for _, ID := range origLabels {
		if labelID == ID {
			continue
		}
		newLabels = append(newLabels, ID)
	}

	bmark.setLabelIDs(newLabels)

	b.add(bmark)
	b.storeBookmarks(userID)

	return nil
}

func (b *Bookmarks) getLabelNames(userID string, bmark *Bookmark) ([]string, error) {
	labels := NewLabels(b.api)
	labels, _ = labels.getLabels(userID)

	var labelNames []string
	for _, id := range bmark.getLabelIDs() {
		name, err := labels.getNameFromID(userID, id)
		if err != nil {
			return nil, err
		}
		labelNames = append(labelNames, name)
	}
	return labelNames, nil
}

func getBookmarksKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreBookmarksKey, userID)
}
