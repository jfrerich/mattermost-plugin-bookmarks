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
func (l *Labels) storeLabels(userID string) error {
	bb, jsonErr := json.Marshal(l)
	if jsonErr != nil {
		return jsonErr
	}

	key := getLabelsKey(userID)
	appErr := l.api.KVSet(key, bb)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

// getNameFromID returns the Name of a Label
func (l *Labels) getNameFromID(userID string, ID string) (string, error) {
	label := l.get(ID)
	return label.Name, nil
}

// getLabels returns a users labels
func (l *Labels) getLabels(userID string) (*Labels, error) {

	// if a user does not have labels, bb will be nil
	bb, appErr := l.api.KVGet(getLabelsKey(userID))
	if appErr != nil {
		return nil, appErr
	}

	labels := NewLabels(l.api)
	if bb == nil {
		return labels, nil
	}

	jsonErr := json.Unmarshal(bb, &labels)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return labels, nil
}

// getLabelByName returns a label with the provided label name
func (l *Labels) getLabelByName(userID string, labelName string) (*Label, error) {

	labels, err := l.getLabels(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	for _, l := range labels.ByID {
		if l.Name == labelName {
			return l, nil
		}
	}

	return nil, nil
}

// getIDsFromNames returns a list of label names
func (l *Labels) getIDsFromNames(userID string, labelNames []string) ([]string, error) {

	newLabelNames := labelNames

	// need to determine which names did not have an ID in the labels store
	// then create them in the store and attach them to the bookmark

	// build array of all UUIDs for the bookmark
	var uuids []string
	for id, label := range l.ByID {
		for _, name := range labelNames {
			if label.Name == name {
				newLabelNames = removeFromArray(label.Name, newLabelNames)
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
			l.add(labelID, label)
			uuids = append(uuids, labelID)
		}
	}

	l.storeLabels(userID)
	return uuids, nil
}

// getIDFromName returns a label name with the corresponding label ID
func (l *Labels) getIDFromName(userID string, labelName string) (string, error) {

	labels, err := l.getLabels(userID)
	if err != nil {
		return "", errors.New(err.Error())
	}

	if labels == nil {
		return "", errors.New(fmt.Sprint("User does not have any labels"))
	}

	// return the labelId if found
	for id, l := range labels.ByID {
		if l.Name == labelName {
			return id, nil
		}
	}
	return "", errors.New(fmt.Sprintf("Label: `%s` does not exist", labelName))
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

// addLabel stores a label into the users label store
func (l *Labels) addLabel(userID string, labelName string) (*Label, error) {

	// check if name already exists
	label, err := l.getLabelByName(userID, labelName)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// User already has label with this labelName
	if label != nil {
		return nil, errors.New(fmt.Sprintf("Label with name `%s` already exists", label.Name))
	}

	// get all labels for user
	labels, err := l.getLabels(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// no labels, initialize the store and save
	if labels == nil {
		labels = NewLabels(l.api) // save first label
	}

	labelID := NewID()
	label = &Label{
		Name: labelName,
	}
	labels.add(labelID, label)

	if err = l.storeLabels(userID); err != nil {
		return nil, errors.New(err.Error())
	}

	return label, nil
}

// deleteByID deletes a label from the store
func (l *Labels) deleteByID(userID, labelID string) error {

	// check if exists
	_, ok := l.exists(labelID)
	if !ok {
		return errors.New(fmt.Sprintf("Label with ID `%s` doesn't exist", labelID))
	}

	l.delete(labelID)
	l.storeLabels(userID)

	return nil
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
