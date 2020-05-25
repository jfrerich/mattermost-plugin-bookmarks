package main

import (
	"regexp"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

type BookmarksFilters struct {
	TitleText  string
	LabelIDs   []string
	LabelNames []string
}

// applyFilters will apply the available filters to an object of bookmarks
func (b *Bookmarks) applyFilters(filters *BookmarksFilters) (*Bookmarks, error) {
	newBmarks := NewBookmarksWithUser(b.api, b.userID)
	// iter through bookmarks
	for _, bmark := range b.ByID {
		filteredBmark := bmark.withLabelIDs(filters.LabelIDs)
		filteredBmark = filteredBmark.withLabelNames(filters.LabelNames, b.api, b.userID)
		filteredBmark = filteredBmark.withTitleText(filters.TitleText)

		if filteredBmark != nil {
			// Do not save the bookmarks to the store. only hold in data structure
			newBmarks.ByID[bmark.PostID] = filteredBmark
		}
	}

	return newBmarks, nil
}

// withLabels returns a bookmark with given label IDs or nil
func (bm *Bookmark) withLabelIDs(ids []string) *Bookmark {
	// return bookmark if no ids requested or bmark is nil
	if ids == nil || bm == nil {
		return bm
	}

	// iter through bmark label ids
	for _, labelID := range bm.getLabelIDs() {
		// iter through requested labelIDs
		for _, id := range ids {
			// return bookmark if has requested labelID
			if labelID == id {
				return bm
			}
		}
	}
	return nil
}

// withLabelNames returns a bookmark with given label names or nil
func (bm *Bookmark) withLabelNames(names []string, api plugin.API, userID string) *Bookmark {
	// return bookmark if no names requested or bmark is nil
	if len(names) == 0 || bm == nil {
		return bm
	}

	labels := NewLabelsWithUser(api, userID)
	labels, _ = labels.getLabels()
	// if err != nil {
	// 	return p.responsef(args, err.Error())
	// }

	// iter through bmark label ids
	for _, labelID := range bm.getLabelIDs() {
		// iter through requested label names
		for _, name := range names {
			// return bookmark if has requested label name
			n, _ := labels.getNameFromID(labelID)
			if n == name {
				return bm
			}
		}
	}
	return nil
}

// withTitleText returns a bookmark with given title text or nil
func (bm *Bookmark) withTitleText(text string) *Bookmark {
	// return bookmark if empty text is empty or bmark is nil
	if text == "" || bm == nil {
		return bm
	}

	title := bm.getTitle()
	re := regexp.MustCompile(text)
	// return bookmark if has requested text
	if re.MatchString(title) {
		return bm
	}

	return nil
}
