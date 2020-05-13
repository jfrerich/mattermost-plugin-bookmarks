package main

import (
	"github.com/mattermost/mattermost-server/v5/plugin"
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

func (l *Labels) add(uuid string, label *Label) {
	l.ByID[uuid] = label
}

func (l *Labels) get(id string) (*Label, error) {
	return l.ByID[id], nil
}

func (l *Labels) delete(id string) {
	delete(l.ByID, id)
}
