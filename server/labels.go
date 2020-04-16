package main

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

const StoreLabelsKey = "labels"

// storeLabels stores all the users labels
func (p *Plugin) storeLabels(userID string, labels *Labels) error {
	jsonBookmarks, jsonErr := json.Marshal(labels)
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

// getLabels returns a users Labels available for all their bookmarks.
func (p *Plugin) getLabels(userID string) (*Labels, error) {

	// if a user does not have bookmarks, bb will be nil
	bb, appErr := p.API.KVGet(getLabelsKey(userID))
	if appErr != nil {
		return nil, appErr
	}

	labels := NewLabels()
	if bb == nil {
		return labels, nil
	}

	jsonErr := json.Unmarshal(bb, &labels)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return labels, nil
}

// addLabels stores labels available for bookmarks
func (p *Plugin) addLabels(userID string, labelNames []string) (*Labels, error) {

	// get all bookmarks for user
	labels, err := p.getLabels(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// no marks, initialize the store first
	if labels == nil {
		labels = NewLabels()
	}

	for _, name := range labelNames {
		label := new(Label)
		label.Name = name
		labels.add(label)
	}

	if err = p.storeLabels(userID, labels); err != nil {
		return nil, errors.New(err.Error())
	}
	return labels, nil
}

func getLabelsKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreLabelsKey, userID)
}
