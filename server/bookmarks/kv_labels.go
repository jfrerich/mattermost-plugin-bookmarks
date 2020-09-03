package bookmarks

import (
	"encoding/json"
	"fmt"
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

// DeleteByID deletes a label from the store
func (l *Labels) DeleteByID(id string) error {
	delete(l.ByID, id)

	if err := l.StoreLabels(); err != nil {
		return err
	}
	return nil
}
