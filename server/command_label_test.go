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

func TestExecuteCommandLabel(t *testing.T) {
	tests := map[string]struct {
		commandArgs       *model.CommandArgs
		bookmarks         *Bookmarks
		expectedMsgPrefix string
		expectedContains  []string
	}{
		"User does not provide label sub-command": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("Missing "),
			expectedContains:  []string{"Missing label sub-command", "bookmarks label add"},
		},

		// ADD
		"ADD User does not provide label names": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label add"},
			bookmarks:         nil,
			expectedMsgPrefix: "",
			// expectedContains:  []string{"Please specify a label name", getHelp(labelCommandText)},
			expectedContains: []string{"Please specify a label name"},
		},
		"User adds first label": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label add label1"},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: "",
			expectedContains:  []string{"Added Label: label1"},
		},
		"User adds 2 labels": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label add label1 label2"},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace("Added Labels:\n"),
			expectedContains:  []string{"Added Labels:", "label1", "label2"},
		},
		// "PostID doesn't exist": {
		// 	commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v", PostIDDoesNotExist)},
		// 	bookmarks:         nil,
		// 	expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("PostID `%v` is not a valid postID", PostIDDoesNotExist)),
		// 	expectedContains:  nil,
		// },

		// VIEW
		"VIEW": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks label view"},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: "#### Labels List",
			expectedContains:  []string{"#### Labels List", "label1"},
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
