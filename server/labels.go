package main

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// getLabels returns a users Labels available for all their bookmarks.
func (p *Plugin) getLabels(userID string) (*Labels, error) {

	// if a user does not have bookmarks, bb will be nil
	bb, appErr := p.API.KVGet(getBookmarksKey(userID))
	if appErr != nil {
		return nil, appErr
	}

	bmarks := NewBookmarks()
	if bb == nil {
		return bmarks.Labels, nil
	}

	jsonErr := json.Unmarshal(bb, &bmarks)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return bmarks.Labels, nil
}

// getLabelsForBookmark returns an array of label names for a given bookmark
func (p *Plugin) getLabelsForBookmark(userID string, bmarkID string) ([]string, error) {

	bmark, err := p.getBookmark(userID, bmarkID)
	if err != nil {
		return nil, err
	}

	return bmark.LabelNames, nil
}

// addLabelsToBookmarks stores labels available for bookmarks
func (p *Plugin) addLabelsToBookmarks(userID string, labels []string) (*Bookmarks, error) {

	// get all bookmarks for user
	bmarks, err := p.getBookmarks(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// no marks, initialize the store first
	if bmarks == nil {
		bmarks = NewBookmarks()
	}

	for _, name := range labels {
		label := new(Label)
		label.Name = name
		bmarks.addLabel(label)
	}

	if err = p.storeBookmarks(userID, bmarks); err != nil {
		return nil, errors.New(err.Error())
	}
	return bmarks, nil
}
