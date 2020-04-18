package main

import (
	"bytes"
	"encoding/base32"
	"encoding/json"
	"fmt"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

// StoreLabelsKey is the key used to store labels in the plugin KV store
const StoreLabelsKey = "labels"

// storeLabels stores all the users labels
func (p *Plugin) storeLabels(userID string, labels *Labels) error {
	bb, jsonErr := json.Marshal(labels)
	if jsonErr != nil {
		return jsonErr
	}

	key := getLabelsKey(userID)
	appErr := p.MattermostPlugin.API.KVSet(key, bb)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

// getLabels returns a users Labels available for all their bookmarks.
func (p *Plugin) getLabelNameByID(userID string, ID string) (string, error) {
	// get all labels for user
	labels, err := p.getLabels(userID)
	if err != nil {
		return "", errors.New(err.Error())
	}

	label := labels.get(ID)
	return label.Name, nil
}

// getLabels returns a users Labels available for all their bookmarks.
func (p *Plugin) getLabels(userID string) (*Labels, error) {

	// if a user does not have labels, bb will be nil
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
func (p *Plugin) getLabelByName(userID string, labelName string) (*Label, error) {

	// get all labels for user
	labels, err := p.getLabels(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if labels == nil {
		return nil, nil
	}

	for _, l := range labels.ByID {
		if l.Name == labelName {
			return l, nil
		}
	}

	return nil, nil
}

// addLabels stores labels available for bookmarks
func (p *Plugin) getLabelIDsFromNames(userID string, labelNames []string) ([]string, error) {

	// // get all labels for user
	labels, err := p.getLabels(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if labels == nil {
		return nil, nil
	}

	newLabelNames := labelNames

	// need to determine which names did not have an ID in the labels store
	// then create them in the store and attach them to the bookmark

	// build array of all UUIDs for the bookmark
	var uuids []string
	for id, l := range labels.ByID {
		for _, name := range labelNames {
			if l.Name == name {
				newLabelNames = removeFromArray(l.Name, newLabelNames)
				uuids = append(uuids, id)
			}
		}
	}

	//generate new labels
	if len(newLabelNames) > 0 {
		for _, name := range newLabelNames {
			labelID := NewID()
			label := &Label{
				Name: name,
			}
			labels.add(labelID, label)
			uuids = append(uuids, labelID)
		}
	}

	p.storeLabels(userID, labels)
	return uuids, nil
}

func removeFromArray(name string, array []string) []string {
	var newArray []string
	for _, elem := range array {
		if name == elem {
			continue
		}
		newArray = append(newArray, elem)
	}

	return newArray
}

// addLabels stores labels available for bookmarks
func (p *Plugin) addLabel(userID string, labelName string) (*Label, error) {

	// check if name already exists
	label, err := p.getLabelByName(userID, labelName)

	// User already has label with this labelName
	if label != nil {
		return nil, errors.New(fmt.Sprintf("Label with name `%s` already exists", label.Name))
	}

	// get all labels for user
	labels, err := p.getLabels(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// no labels, initialize the store and save
	if labels == nil {
		labels = NewLabels() // save first label
	}

	labelID := NewID()
	label = &Label{
		Name: labelName,
	}
	labels.add(labelID, label)

	if err = p.storeLabels(userID, labels); err != nil {
		return nil, errors.New(err.Error())
	}

	return label, nil
}

// deleteLabelByName deletes a label from the store
func (p *Plugin) deleteLabelByName(userID, labelName string) (*Label, error) {

	labels, err := p.getLabels(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if labels == nil {
		return nil, errors.New(fmt.Sprintf("User doesn't have any labels"))
	}

	// check if exists
	label, err := p.getLabelByName(userID, labelName)
	if label == nil {
		return nil, errors.New(fmt.Sprintf("Label with name `%s` doesn't exist", labelName))
	}

	labelID, _ := p.getLabelIDsFromNames(userID, []string{labelName})
	labels.delete(labelID[0])
	p.storeLabels(userID, labels)

	return label, nil
}

func getLabelsKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreLabelsKey, userID)
}

var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769")

// NewID is a globally unique identifier.  It is a [A-Z0-9] string 26
// characters long.  It is a UUID version 4 Guid that is zbased32 encoded
// with the padding stripped off.
func NewID() string {
	var b bytes.Buffer
	encoder := base32.NewEncoder(encoding, &b)
	encoder.Write(uuid.NewRandom())
	encoder.Close()
	b.Truncate(26) // removes the '==' padding
	return b.String()
}
