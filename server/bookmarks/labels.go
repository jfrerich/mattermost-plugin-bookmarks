package bookmarks

import (
	"encoding/json"
	"fmt"

	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/utils"
	"github.com/pkg/errors"
)

// Labels contains a map of labels with the label name as the key
type Labels struct {
	ByID   map[string]*Label
	api    pluginapi.API
	userID string
}

// Label defines the parameters of a label
type Label struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	// Color string `json:"color"`
}

// NewLabels returns an initialized Labels struct
func NewLabels(userID string) *Labels {
	return &Labels{
		ByID:   make(map[string]*Label),
		userID: userID,
	}
}

// NewLabelsWithUser returns an initialized Labels for a User
func NewLabelsWithUser(api pluginapi.API, userID string) (*Labels, error) {
	bb, appErr := api.KVGet(GetLabelsKey(userID))
	if appErr != nil {
		return nil, errors.Wrapf(appErr, "Unable to get labels for user %s", userID)
	}

	userLabels, err := LabelsFromJSON(bb)
	if err != nil {
		return nil, err
	}
	userLabels.api = api
	userLabels.userID = userID

	return userLabels, nil
}

// LabelsFromJSON returns unmarshalled bookmark or initialized bookmarks if
// bytes are empty
func LabelsFromJSON(bytes []byte) (*Labels, error) {
	labels := &Labels{
		ByID: make(map[string]*Label),
	}

	if len(bytes) != 0 {
		jsonErr := json.Unmarshal(bytes, &labels)
		if jsonErr != nil {
			return nil, jsonErr
		}
	}
	return labels, nil
}

// GetNameFromID returns the Name of a Label
func (l *Labels) GetNameFromID(id string) (string, error) {
	label, err := l.Get(id)
	if err != nil {
		return "", err
	}
	if label == nil {
		return "", nil
	}

	return label.Name, nil
}

// GetLabelByName returns a label with the provided label name
func (l *Labels) GetLabelByName(labelName string) *Label {
	if l == nil {
		return nil
	}
	for _, label := range l.ByID {
		if label.Name == labelName {
			return label
		}
	}
	return nil
}

// GetIDFromName returns a label name with the corresponding label ID
func (l *Labels) GetIDFromName(labelName string) (string, error) {
	if l == nil {
		return "", errors.New("user does not have any labels")
	}

	// return the labelId if found
	for id, label := range l.ByID {
		if label.Name == labelName {
			return id, nil
		}
	}
	return "", errors.New(fmt.Sprintf("Label: `%s` does not exist", labelName))
}

// addLabel stores a label into the users label store
func (l *Labels) AddLabel(labelName string) (*Label, error) {
	// check if name already exists
	label := l.GetLabelByName(labelName)

	// User already has label with this labelName
	if label != nil {
		return nil, errors.New(fmt.Sprintf("Label with name `%s` already exists", label.Name))
	}

	labelID := utils.NewID()
	label = &Label{
		Name: labelName,
		ID:   labelID,
	}
	if err := l.Add(labelID, label); err != nil {
		return nil, err
	}

	return label, nil
}

// DeleteByID deletes a label from the store
func (l *Labels) DeleteByID(labelID string) error {
	if err := l.delete(labelID); err != nil {
		return err
	}
	return nil
}
