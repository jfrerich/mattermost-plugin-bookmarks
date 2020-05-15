package main

import (
	"encoding/json"

	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

// Labels contains a map of labels with the label name as the key
type Labels struct {
	ByID   map[string]*Label
	api    plugin.API
	userID string
}

// Label defines the parameters of a label
type Label struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	// Color string `json:"color"`
}

// NewLabels returns an initialized Labels struct
func NewLabels(api plugin.API) *Labels {
	return &Labels{
		ByID: make(map[string]*Label),
		api:  api,
	}
}

// NewLabelsWithUser returns an initialized Labels for a User
func NewLabelsWithUser(api plugin.API, userID string) *Labels {
	return &Labels{
		ByID:   make(map[string]*Label),
		api:    api,
		userID: userID,
	}
}

func (l *Labels) add(uuid string, label *Label) error {
	l.ByID[uuid] = label
	if err := l.storeLabels(); err != nil {
		return errors.Wrap(err, "failed to add label")
	}
	return nil
}

func (l *Labels) get(id string) (*Label, error) {
	return l.ByID[id], nil
}

func (l *Labels) delete(id string) error {
	delete(l.ByID, id)
	if err := l.storeLabels(); err != nil {
		return err
	}
	return nil
}

// storeLabels stores all the users labels
func (l *Labels) storeLabels() error {
	bb, jsonErr := json.Marshal(l)
	if jsonErr != nil {
		return jsonErr
	}

	key := getLabelsKey(l.userID)
	appErr := l.api.KVSet(key, bb)
	if appErr != nil {
		return appErr
	}

	return nil
}
