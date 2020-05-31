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

// getBookmark returns a bookmark with the specified bookmarkID
func (b *Bookmarks) getBookmark(bmarkID string) (*Bookmark, error) {
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
func (b *Bookmarks) addBookmark(bmark *Bookmark) error {
	// user doesn't have any bookmarks add first bookmark and return
	if len(b.ByID) == 0 {
		if err := b.add(bmark); err != nil {
			return err
		}
		return nil
	}

	// bookmark already exists, update ModifiedAt and save
	_, ok := b.exists(bmark.PostID)
	if ok {
		b.updateTimes(bmark.PostID)
		b.updateLabels(bmark)
		if err := b.add(bmark); err != nil {
			return err
		}
		return nil
	}

	// bookmark doesn't exist. Add it
	if err := b.add(bmark); err != nil {
		return err
	}
	return nil
}

// BookmarksFromJson returns unmarshalled bookmark or initialized bookmarks if
// bytes are emtpy
func (b *Bookmarks) BookmarksFromJson(bytes []byte) (*Bookmarks, error) {
	bmarks := NewBookmarksWithUser(b.api, b.userID)
	if len(bytes) != 0 {
		jsonErr := json.Unmarshal(bytes, &bmarks)
		if jsonErr != nil {
			return nil, jsonErr
		}
	}
	return bmarks, nil
}

// getBookmarks returns a users bookmarks.  If the user has no bookmarks,
// return nil bookmarks
func (b *Bookmarks) getBookmarks() (*Bookmarks, error) {
	// if a user not not have bookmarks, bb will be nil
	bb, appErr := b.api.KVGet(getBookmarksKey(b.userID))
	if appErr != nil {
		return nil, errors.Wrapf(appErr, "Unable to get bookmarks for user %s", b.userID)
	}

	return b.BookmarksFromJson(bb)
}

// ByPostCreateAt returns an array of bookmarks sorted by post.CreateAt times
func (b *Bookmarks) ByPostCreateAt() ([]*Bookmark, error) {
	// build temp map
	tempMap := make(map[int64]string)
	for _, bmark := range b.ByID {
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
		bmark := b.ByID[tempMap[int64(k)]]
		bookmarks = append(bookmarks, bmark)
	}

	return bookmarks, nil
}

func (b *Bookmarks) getBookmarksWithLabelID(labelID string) (*Bookmarks, error) {
	bmarksWithLabel := NewBookmarksWithUser(b.api, b.userID)

	for _, bmark := range b.ByID {
		if bmark.hasLabels(bmark) {
			for _, id := range bmark.getLabelIDs() {
				if labelID == id {
					if err := bmarksWithLabel.add(bmark); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return bmarksWithLabel, nil
}

// deleteBookmark deletes a bookmark from the store
func (b *Bookmarks) deleteBookmark(bmarkID string) (*Bookmark, error) {
	var bmark *Bookmark

	_, ok := b.exists(bmarkID)
	if !ok {
		return bmark, errors.New(fmt.Sprintf("Bookmark `%v` does not exist", bmarkID))
	}

	bmark, err := b.getBookmark(bmarkID)
	if err != nil {
		return nil, err
	}

	b.delete(bmarkID)
	if err := b.storeBookmarks(); err != nil {
		return nil, err
	}

	return bmark, nil
}

// deleteLabel deletes a label from a bookmark
func (b *Bookmarks) deleteLabel(bmarkID string, labelID string) error {
	bmark, err := b.getBookmark(bmarkID)
	if err != nil {
		return err
	}

	origLabels := bmark.getLabelIDs()

	var newLabels []string
	for _, ID := range origLabels {
		if labelID == ID {
			continue
		}
		newLabels = append(newLabels, ID)
	}

	bmark.addLabelIDs(newLabels)

	if err := b.add(bmark); err != nil {
		return err
	}

	return nil
}

// getLabelNames returns an array of labelNames for a given bookmark
func (b *Bookmarks) getBmarkLabelNames(bmark *Bookmark) ([]string, error) {
	l, err := NewLabelsWithUser(b.api, b.userID).getLabels()
	if err != nil {
		return nil, err
	}

	var labelNames []string
	for _, id := range bmark.getLabelIDs() {
		name, err := l.getNameFromID(id)
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
