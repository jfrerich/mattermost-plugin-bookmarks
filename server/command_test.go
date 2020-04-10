package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExecuteCommandView(t *testing.T) {
	tests := map[string]struct {
		commandArgs       *model.CommandArgs
		bookmarks         *Bookmarks
		expectedMsgPrefix string
		expectedContains  []string
	}{

		// /bookmarks view testing
		"VIEW User has 3 bookmarks": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view"},
			bookmarks:         getTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace("#### Bookmarks List"),
			expectedContains:  []string{"Bookmarks List", "ID1", "ID2", "ID3"},
		},
		"VIEW User has no bookmarks": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("You do not have any saved bookmarks"),
			expectedContains:  nil,
		},

		// /bookmarks remove testing
		"REMOVE User tries to delete a bookmark but has none": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks remove bmarkID"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("User doesn't have any bookmarks"),
			expectedContains:  nil,
		},
		"REMOVE User has bmarks tries to delete bookmark that doesnt exist": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks remove bmarkID"},
			bookmarks:         getTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace("Bookmark `bmarkID` does not exist"),
			expectedContains:  nil,
		},
		"REMOVE User successfully deletes 1 bookmark": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks remove ID2"},
			bookmarks:         getTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace("Removed bookmark ID: ID2"),
			expectedContains:  nil,
		},
	}
	for name, tt := range tests {
		api := makeAPIMock()
		tt.commandArgs.UserId = "junkID"
		siteURL := "https://myhost.com"
		teamID1 := "teamID1"
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
					for i, _ := range tt.expectedContains {
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
