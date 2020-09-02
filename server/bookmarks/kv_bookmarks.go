package bookmarks

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// StoreBookmarksKey is the key used to store bookmarks in the plugin KV store
const StoreBookmarksKey = "bookmarks"

func GetBookmarksKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreBookmarksKey, userID)
}

// storeBookmarks stores all the users bookmarks
func (b *Bookmarks) StoreBookmarks() error {
	bb, jsonErr := json.Marshal(b)
	if jsonErr != nil {
		return jsonErr
	}

	key := GetBookmarksKey(b.userID)
	appErr := b.api.KVSet(key, bb)
	if appErr != nil {
		return appErr
	}

	return nil
}

// BookmarksFromJSON returns unmarshalled bookmark or initialized bookmarks if
// bytes are empty
func BookmarksFromJSON(bytes []byte) (*Bookmarks, error) {
	bmarks := &Bookmarks{
		ByID: make(map[string]*Bookmark),
	}

	if len(bytes) != 0 {
		jsonErr := json.Unmarshal(bytes, &bmarks)
		if jsonErr != nil {
			return nil, jsonErr
		}
	}
	return bmarks, nil
}

// DeleteBookmark deletes a bookmark from the store
func (b *Bookmarks) DeleteBookmark(bmarkID string) (*Bookmark, error) {
	var bmark *Bookmark

	_, ok := b.exists(bmarkID)
	if !ok {
		return bmark, errors.New(fmt.Sprintf("Bookmark `%v` does not exist", bmarkID))
	}

	bmark, err := b.GetBookmark(bmarkID)
	if err != nil {
		return nil, err
	}

	b.delete(bmarkID)
	if err := b.StoreBookmarks(); err != nil {
		return nil, err
	}

	return bmark, nil
}

func (b *Bookmarks) Add(bmark *Bookmark) error {
	b.ByID[bmark.PostID] = bmark

	if err := b.StoreBookmarks(); err != nil {
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