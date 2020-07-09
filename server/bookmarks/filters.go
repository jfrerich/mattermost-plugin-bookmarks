package bookmarks

import (
	"regexp"

	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi"
)

type Filters struct {
	TitleText  string
	LabelIDs   []string
	LabelNames []string
}

// ApplyFilters will apply the available filters to an object of bookmarks
func (b *Bookmarks) ApplyFilters(filters *Filters) (*Bookmarks, error) {
	// FIXME: This should not require setting the api again.
	bmarks := NewBookmarks(b.userID)
	bmarks.api = b.api

	// iter through bookmarks
	for _, bmark := range b.ByID {
		filteredBmark := bmark.withLabelIDs(filters.LabelIDs)
		filteredBmark = filteredBmark.withLabelNames(filters.LabelNames, b.api, b.userID)
		filteredBmark = filteredBmark.withTitleText(filters.TitleText)

		if filteredBmark != nil {
			// Do not save the bookmarks to the store. only hold in data structure
			bmarks.ByID[bmark.PostID] = filteredBmark
		}
	}

	return bmarks, nil
}

// withLabels returns a bookmark with given label IDs or nil
func (bm *Bookmark) withLabelIDs(ids []string) *Bookmark {
	// return bookmark if no ids requested or bmark is nil
	if ids == nil || bm == nil {
		return bm
	}

	// iter through bmark label ids
	for _, labelID := range bm.GetLabelIDs() {
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
func (bm *Bookmark) withLabelNames(names []string, api pluginapi.API, userID string) *Bookmark {
	// return bookmark if no names requested or bmark is nil
	if len(names) == 0 || bm == nil {
		return bm
	}

	labels, _ := NewLabelsWithUser(api, userID)
	// if err != nil {
	// 	return p.responsef(args, err.Error())
	// }

	// iter through bmark label ids
	for _, labelID := range bm.GetLabelIDs() {
		// iter through requested label names
		for _, name := range names {
			// return bookmark if has requested label name
			n, _ := labels.GetNameFromID(labelID)
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

	title := bm.GetTitle()
	re := regexp.MustCompile(text)
	// return bookmark if has requested text
	if re.MatchString(title) {
		return bm
	}

	return nil
}
