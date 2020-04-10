package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	PostIDDoesNotExist = "PostIDDoesNotExist"
	PostIDExists       = "ID2"
	UserID             = "UserID"
	teamID1            = "teamID1"
)

func TestExecuteCommandView(t *testing.T) {
	tests := map[string]struct {
		commandArgs       *model.CommandArgs
		bookmarks         *Bookmarks
		expectedMsgPrefix string
		expectedContains  []string
	}{
		// Unknown Slash Commmand
		"UNKNOWN slash commnad": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks UnknownCommand"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("Unknown command: /bookmarks UnknownCommand"),
			expectedContains:  []string{},
		},

		// Help Slash Commmand
		"HELP slash commnad": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks help"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("###### Bookmarks Slash Command Help "),
			expectedContains:  []string{"bookmarks add", "bookmarks view", "bookmarks remove"},
		},

		// ADD Slash Commmand
		"ADD User doesn't provide an ID": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks add"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("Missing "),
			expectedContains:  []string{"Missing sub-command", "bookmarks add"},
		},
		"ADD PostID doesn't exist": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v", PostIDDoesNotExist)},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("PostID `%v` is not a valid postID", PostIDDoesNotExist)),
			expectedContains:  nil,
		},
		"ADD Boommark added": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v", PostIDExists)},
			bookmarks:         getTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("Added bookmark: {PostID:%v Title:This message exists", PostIDExists)),
			expectedContains:  nil,
		},

		// VIEW Slash Command
		"VIEW User has no bookmarks": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("You do not have any saved bookmarks"),
			expectedContains:  nil,
		},
		"VIEW User has 3 bookmarks": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view"},
			bookmarks:         getTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace("#### Bookmarks List"),
			expectedContains:  []string{"Bookmarks List", "ID1", "ID2", "ID3"},
		},

		// REMOVE Slash Command
		"REMOVE User doesn't provide an ID": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks remove"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("Missing "),
			expectedContains:  []string{"Missing sub-command", "bookmarks remove"},
		},
		"REMOVE User tries to delete a bookmark but has none": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks remove bmarkID"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("User doesn't have any bookmarks"),
			expectedContains:  nil,
		},
		"REMOVE User has bmarks tries to delete bookmark that doesnt exist": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks remove %v", PostIDDoesNotExist)},
			bookmarks:         getTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("Bookmark `%v` does not exist", PostIDDoesNotExist)),
			expectedContains:  nil,
		},
		"REMOVE User successfully deletes 1 bookmark": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks remove %v", PostIDExists)},
			bookmarks:         getTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("Removed bookmark ID: %v", PostIDExists)),
			expectedContains:  nil,
		},
	}
	for name, tt := range tests {
		api := makeAPIMock()
		tt.commandArgs.UserId = UserID
		siteURL := "https://myhost.com"
		api.On("GetPost", PostIDDoesNotExist).Return(nil, &model.AppError{Message: "An Error Occurred"})
		api.On("GetPost", PostIDExists).Return(&model.Post{Message: "This message exists"}, nil)
		api.On("addBookmark", UserID, tt.bookmarks).Return(mock.Anything)
		api.On("GetTeam", mock.Anything).Return(&model.Team{Id: teamID1}, nil)
		api.On("GetConfig", mock.Anything).Return(&model.Config{ServiceSettings: model.ServiceSettings{SiteURL: &siteURL}})
		api.On("exists", mock.Anything).Return(true)
		// api.On("ByID", mock.Anything).Return(true)

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

func makeAPIMock() *plugintest.API {
	api := &plugintest.API{}

	api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
	api.On("LogWarn", mock.Anything, mock.Anything, mock.Anything).Maybe()
	api.On("LogError", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

	return api
}
