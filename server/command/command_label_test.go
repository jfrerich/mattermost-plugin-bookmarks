package command

import (
	"encoding/json"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/bookmarks"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi/mock_pluginapi"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/assert"
)

func TestExecuteCommandLabel(t *testing.T) {
	tests := map[string]struct {
		command             string
		bookmarks           *bookmarks.Bookmarks
		labels              *bookmarks.Labels
		expectedMsgPrefix   string
		expectedContains    []string
		expectedNotContains []string
	}{
		"User does not provide label sub-command": {
			command:           "/bookmarks label",
			labels:            nil,
			expectedMsgPrefix: strings.TrimSpace("Missing "),
			expectedContains:  []string{"Missing label sub-command", "bookmarks label add"},
		},

		// ADD
		"ADD User does not provide label names": {
			command:           "/bookmarks label add",
			labels:            nil,
			expectedMsgPrefix: "",
			expectedContains:  []string{"Please specify a label name"},
		},
		"ADD User adds first label": {
			command:           "/bookmarks label add label9",
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "",
			expectedContains:  []string{"Added Label: label9"},
		},
		"ADD User tries creating label with name that already exists": {
			command:           "/bookmarks label add label1",
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "",
			expectedContains:  []string{"Label with name `label1` already exists"},
		},
		"ADD User adds one label successfully with existing labels": {
			command:           "/bookmarks label add NewLabelName",
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "",
			expectedContains:  []string{"Added Label: NewLabelName"},
		},

		// RENAME - successfully renames a label
		"RENAME User provides only one label": {
			command:           "/bookmarks label rename label1",
			labels:            nil,
			expectedMsgPrefix: "Please specify a `to` and `from` label name",
			expectedContains:  nil,
		},
		"RENAME User tries renaming label that doesn't exist": {
			command:           "/bookmarks label rename label1 label2",
			labels:            &bookmarks.Labels{},
			expectedMsgPrefix: "Label `label1` does not exist",
			expectedContains:  nil,
		},
		"RENAME User tries renaming a label to itself": {
			command:           "/bookmarks label rename label1 label1",
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "Cannot rename Label `label1` to `label1`. Label already exists. Please choose a different label name",
			expectedContains:  nil,
		},
		"RENAME User tries renaming a label to another label name that exists": {
			command:           "/bookmarks label rename label1 label2",
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "Cannot rename Label `label1` to `label2`. Label already exists. Please choose a different label name",
			expectedContains:  nil,
		},
		"RENAME User successfully renames label": {
			command:           "/bookmarks label rename label1 labelthatdoesntexist",
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "Renamed label from `label1` to `labelthatdoesntexist`",
			expectedContains:  nil,
		},

		// REMOVE - user does not have any saved labels
		"REMOVE User does not provide label name": {
			command:           "/bookmarks label remove",
			labels:            nil,
			expectedMsgPrefix: "",
			expectedContains:  []string{"Please specify a label name"},
		},
		"REMOVE User tries to remove a label but has none": {
			command:           "/bookmarks label remove JunkLabel",
			bookmarks:         &bookmarks.Bookmarks{},
			labels:            &bookmarks.Labels{},
			expectedMsgPrefix: "You do not have any saved labels",
			expectedContains:  nil,
		},

		// REMOVE - user has saved labels
		"REMOVE User tries to remove a label that does not exist": {
			command:           "/bookmarks label remove labeldoesnotexist",
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "",
			expectedContains:  []string{"Label: `labeldoesnotexist` does not exist"},
		},
		"REMOVE User successfully removes a label that exists": {
			command:           "/bookmarks label remove label1",
			bookmarks:         &bookmarks.Bookmarks{},
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "",
			expectedContains:  []string{"Removed label: `label1`"},
		},
		"REMOVE User tries to remove a label that exists in a bookmark": {
			command:           "/bookmarks label remove label1",
			labels:            getExecuteCommandTestLabels(),
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: "There are 2 bookmarks with the label: `label1`. Use the `--force` flag remove the label from the bookmarks.",
			expectedContains:  nil,
		},
		"REMOVE User tries to remove a label that exists in a bookmark using the force flag": {
			command:           "/bookmarks label remove label1 --force",
			labels:            getExecuteCommandTestLabels(),
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: "Removed label: `label1`",
			expectedContains:  nil,
		},

		// VIEW
		"VIEW User doesn't have any labels": {
			command:           "/bookmarks label view",
			labels:            &bookmarks.Labels{},
			expectedMsgPrefix: "You do not have any saved labels",
			expectedContains:  nil,
		},
		"VIEW User has 2 label": {
			command:           "/bookmarks label view",
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "",
			expectedContains:  []string{"#### Labels List", "label1", "label2"},
		},
	}
	for name, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPluginAPI := mock_pluginapi.NewMockAPI(ctrl)

		jsonBmarks, err := json.Marshal(tt.bookmarks)
		mockPluginAPI.EXPECT().KVGet(bookmarks.GetBookmarksKey(UserID)).Return(jsonBmarks, nil).AnyTimes()

		jsonLabels, err := json.Marshal(tt.labels)
		mockPluginAPI.EXPECT().KVGet(bookmarks.GetLabelsKey(UserID)).Return(jsonLabels, nil).AnyTimes()
		// api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		mockPluginAPI.EXPECT().KVSet(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		t.Run(name, func(t *testing.T) {
			assert.Nil(t, err)
			testCommand := Command{
				Args: &model.CommandArgs{
					UserId:  UserID,
					Command: tt.command},
				API: mockPluginAPI,
			}

			// just check output message.  We don't need to run p.ExecuteCommand()
			message := testCommand.Handle()
			actual := strings.TrimSpace(message)
			assert.True(t, strings.HasPrefix(actual, tt.expectedMsgPrefix), "Expected returned message to start with: \n%s\nActual:\n%s", tt.expectedMsgPrefix, actual)

			if tt.expectedNotContains != nil {
				for i := range tt.expectedNotContains {
					assert.NotContains(t, actual, tt.expectedNotContains[i])
				}
			}
			assert.Contains(t, actual, tt.expectedMsgPrefix)

			if tt.expectedContains != nil {
				for i := range tt.expectedContains {
					assert.Contains(t, actual, tt.expectedContains[i])
				}
			}
		})
	}
}
