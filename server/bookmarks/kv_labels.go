package bookmarks

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// StoreLabelsKey is the key used to store labels in the plugin KV store
const StoreLabelsKey = "labels"

func GetLabelsKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreLabelsKey, userID)
}

// storeLabels stores all the users labels
func (l *Labels) StoreLabels() error {
	bb, jsonErr := json.Marshal(l)
	if jsonErr != nil {
		return jsonErr
	}

	key := GetLabelsKey(l.userID)
	appErr := l.api.KVSet(key, bb)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (l *Labels) Add(uuid string, label *Label) error {
	l.ByID[uuid] = label

	if err := l.StoreLabels(); err != nil {
		return errors.Wrap(err, "failed to add label")
	}
	return nil
}

func (l *Labels) Get(id string) (*Label, error) {
	if l == nil {
		return nil, nil
	}
	return l.ByID[id], nil
}

func (l *Labels) delete(id string) error {
	delete(l.ByID, id)

	if err := l.StoreLabels(); err != nil {
		return err
	}
	return nil
}
