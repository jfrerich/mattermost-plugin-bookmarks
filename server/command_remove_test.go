package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExecuteCommandRemove(t *testing.T) {
	tests := map[string]struct {
		commandArgs       *model.CommandArgs
		bookmarks         *Bookmarks
		expectedMsgPrefix string
		expectedContains  []string
	}{
		"User doesn't provide an ID": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks remove"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("Missing "),
			expectedContains:  []string{"Missing sub-command", "bookmarks remove"},
		},
		"User tries to delete a bookmark but has none": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks remove bmarkID"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("User doesn't have any bookmarks"),
			expectedContains:  nil,
		},
		"User has bmarks tries to delete bookmark that doesnt exist": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks remove %v", PostIDDoesNotExist)},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("Bookmark `%v` does not exist", PostIDDoesNotExist)),
			expectedContains:  nil,
		},
		"User successfully deletes 1 bookmark": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks remove %v", PostIDExists)},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace("Removed bookmark: [:link:](https://myhost.com/_redirect/pl/ID2) `label1` `label2` **_Title2 - "),
			expectedContains:  nil,
		},
		"User successfully deletes 3 bookmark": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks remove %v %v %v", p1ID, p2ID, p3ID)},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: "",
			expectedContains: []string{
				"Removed bookmarks:",
				"[:link:](https://myhost.com/_redirect/pl/ID1) `label1` `label2` **_Title1 - New Bookmark - times are zero",
				"[:link:](https://myhost.com/_redirect/pl/ID2) `label1` `label2` **_Title2 - bookmarks initialized. Times created and same",
				"[:link:](https://myhost.com/_redirect/pl/ID3) **_Title3 - bookmarks already updated once_**",
			},
		},
	}
	for name, tt := range tests {
		api := makeAPIMock()
		tt.commandArgs.UserId = UserID
		siteURL := "https://myhost.com"
		api.On("GetPost", PostIDDoesNotExist).Return(nil, &model.AppError{Message: "An Error Occurred"})
		api.On("GetPost", p1ID).Return(&model.Post{Message: "this is the post.Message"}, nil)
		api.On("GetPost", p2ID).Return(&model.Post{Message: "this is the post.Message"}, nil)
		api.On("GetPost", p3ID).Return(&model.Post{Message: "this is the post.Message"}, nil)
		api.On("GetPost", p4ID).Return(&model.Post{Message: "this is the post.message"}, nil)
		api.On("addBookmark", UserID, tt.bookmarks).Return(mock.Anything)
		api.On("GetTeam", mock.Anything).Return(&model.Team{Id: teamID1}, nil)
		api.On("GetConfig", mock.Anything).Return(&model.Config{ServiceSettings: model.ServiceSettings{SiteURL: &siteURL}})
		api.On("exists", mock.Anything).Return(true)

		jsonBmarks, err := json.Marshal(tt.bookmarks)
		api.On("KVGet", getBookmarksKey(tt.commandArgs.UserId)).Return(jsonBmarks, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)

		labels := getExecuteCommandTestLabels()
		jsonLabels, err := json.Marshal(labels)
		api.On("KVGet", getLabelsKey(tt.commandArgs.UserId)).Return(jsonLabels, nil)

		t.Run(name, func(t *testing.T) {
			assert.Nil(t, err)
			api.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.AnythingOfType("*model.Post")).Run(func(args mock.Arguments) {
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

			p := makePlugin(api)
			cmdResponse, appError := p.ExecuteCommand(&plugin.Context{}, tt.commandArgs)
			require.Nil(t, appError)
			require.NotNil(t, cmdResponse)
		})
	}
}
