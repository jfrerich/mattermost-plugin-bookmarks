package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func getExecuteCommandTestLabels() *Labels {
	l1 := &Label{
		Name: "label1",
	}
	l2 := &Label{
		Name: "label2",
	}
	l3 := &Label{
		Name: "label8",
	}

	api := makeAPIMock()
	labels := NewLabels(api)
	labels.add("UUID1", l1)
	labels.add("UUID2", l2)
	labels.add("UUID3", l3)

	return labels
}

func TestExecuteCommandLabel(t *testing.T) {
	tests := map[string]struct {
		commandArgs       *model.CommandArgs
		bookmarks         *Bookmarks
		labels            *Labels
		expectedMsgPrefix string
		expectedContains  []string
	}{
		"User does not provide label sub-command": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label"},
			labels:            nil,
			expectedMsgPrefix: strings.TrimSpace("Missing "),
			expectedContains:  []string{"Missing label sub-command", "bookmarks label add"},
		},

		// ADD
		"ADD User does not provide label names": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label add"},
			labels:            nil,
			expectedMsgPrefix: "",
			expectedContains:  []string{"Please specify a label name"},
		},
		"ADD User adds first label": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label add label9"},
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "",
			expectedContains:  []string{"Added Label: label9"},
		},
		"ADD User tries creating label with name that already exists": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label add label1"},
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "",
			expectedContains:  []string{"Label with name `label1` already exists"},
		},
		"ADD User adds one label successfully with existing labels": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label add NewLabelName"},
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "",
			expectedContains:  []string{"Added Label: NewLabelName"},
		},

		// REMOVE - user does not have any saved labels
		"REMOVE User does not provide label name": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label remove"},
			labels:            nil,
			expectedMsgPrefix: "",
			expectedContains:  []string{"Please specify a label name"},
		},
		"REMOVE User tries to remove a label but has none": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label remove JunkLabel"},
			bookmarks:         nil,
			labels:            nil,
			expectedMsgPrefix: "You do not have any saved labels",
			expectedContains:  nil,
		},

		// REMOVE - user has saved labels
		"REMOVE User tries to remove a label that does not exist": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label remove labeldoesnotexist"},
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "",
			expectedContains:  []string{"Label: `labeldoesnotexist` does not exist"},
		},
		"REMOVE User successfully removes a label that exists": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label remove label1"},
			bookmarks:         nil,
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "",
			expectedContains:  []string{"Removed label: `label1`"},
		},
		"REMOVE User tries to remove a label that exists in a bookmark": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label remove label1"},
			labels:            getExecuteCommandTestLabels(),
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: "There are 2 bookmarks with the label: `label1`. Use the --force flag remove the label from the bookmarks.",
			expectedContains:  nil,
		},
		"REMOVE User tries to remove a label that exists in a bookmark using the force flag": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label remove label1 --force"},
			labels:            getExecuteCommandTestLabels(),
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: "Removed label: `label1`",
			expectedContains:  nil,
		},

		// VIEW
		"VIEW User doesn't have any labels": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label view"},
			labels:            nil,
			expectedMsgPrefix: "You do not have any saved labels",
			expectedContains:  nil,
		},
		"VIEW User has 2 label": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label view"},
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: "",
			expectedContains:  []string{"#### Labels List", "label1", "label2"},
		},
	}
	for name, tt := range tests {
		api := makeAPIMock()
		tt.commandArgs.UserId = UserID
		siteURL := "https://myhost.com"
		api.On("GetConfig", mock.Anything).Return(&model.Config{ServiceSettings: model.ServiceSettings{SiteURL: &siteURL}})
		api.On("exists", mock.Anything).Return(true)

		jsonBmarks, err := json.Marshal(tt.bookmarks)
		api.On("KVGet", getBookmarksKey(tt.commandArgs.UserId)).Return(jsonBmarks, nil)

		bb, err := json.Marshal(tt.labels)
		api.On("KVGet", getLabelsKey(tt.commandArgs.UserId)).Return(bb, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)

		t.Run(name, func(t *testing.T) {
			assert.Nil(t, err)
			// isSendEphemeralPostCalled := false
			api.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.AnythingOfType("*model.Post")).Run(func(args mock.Arguments) {
				// isSendEphemeralPostCalled = true

				post := args.Get(1).(*model.Post)
				actual := strings.TrimSpace(post.Message)
				assert.True(t, strings.HasPrefix(actual, tt.expectedMsgPrefix), "Expected returned message to start with: \n%s\nActual:\n%s", tt.expectedMsgPrefix, actual)
				if tt.expectedContains != nil {
					for i := range tt.expectedContains {
						assert.Contains(t, actual, tt.expectedContains[i])
					}
				}
				// assert.Contains(t, actual, tt.expectedMsgPrefix)
			}).Once().Return(&model.Post{})
			// assert.Equal(t, true, isSendEphemeralPostCalled)

			p := makePlugin(api)
			cmdResponse, appError := p.ExecuteCommand(&plugin.Context{}, tt.commandArgs)
			require.Nil(t, appError)
			require.NotNil(t, cmdResponse)
			// assert.True(t, isSendEphemeralPostCalled)
		})
	}
}
