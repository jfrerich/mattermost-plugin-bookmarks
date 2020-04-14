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

	b1ID = "ID1"
	b2ID = "ID2"
	b3ID = "ID3"
	b4ID = "ID4"

	b1Title = "Title1 - New Bookmark - times are zero"
	b2Title = "Title2 - bookmarks initialized. Times created and same"
	b3Title = "Title3 - bookmarks already updated once"

	p4Title = "This is a message from a post"
)

func getExecuteCommandTestBookmarks() *Bookmarks {
	bmarks := NewBookmarks()

	b1 := &Bookmark{
		PostID:     b1ID,
		Title:      b1Title,
		LabelNames: []string{"label1", "label2"},
	}
	b2 := &Bookmark{
		PostID:     b2ID,
		Title:      b2Title,
		CreateAt:   model.GetMillis(),
		ModifiedAt: model.GetMillis(),
		LabelNames: []string{"label1", "label2"},
	}
	b3 := &Bookmark{
		PostID:     b3ID,
		Title:      b3Title,
		CreateAt:   model.GetMillis(),
		ModifiedAt: model.GetMillis(),
	}
	b4 := &Bookmark{
		PostID:     b4ID,
		CreateAt:   model.GetMillis(),
		ModifiedAt: model.GetMillis(),
	}

	bmarks.add(b1)
	bmarks.add(b2)
	bmarks.add(b3)
	bmarks.add(b4)

	l1 := &Label{
		Name: "label1",
	}

	bmarks.Labels.add(l1)

	return bmarks
}

func TestExecuteCommand(t *testing.T) {
	tests := map[string]struct {
		commandArgs       *model.CommandArgs
		bookmarks         *Bookmarks
		expectedMsgPrefix string
		expectedContains  []string
	}{
		// No Slash Commmand
		"No slash command": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("###### Bookmarks Slash Command Help "),
			expectedContains:  []string{"bookmarks add", "bookmarks view", "bookmarks remove"},
		},

		// Unknown Slash Commmand
		"UNKNOWN slash command": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks UnknownCommand"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("Unknown command: /bookmarks UnknownCommand"),
			expectedContains:  []string{},
		},

		// Help Slash Commmand
		"HELP slash command": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks help"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("###### Bookmarks Slash Command Help "),
			expectedContains:  []string{"bookmarks add", "bookmarks view", "bookmarks remove"},
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
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("Bookmark `%v` does not exist", PostIDDoesNotExist)),
			expectedContains:  nil,
		},
		"REMOVE User successfully deletes 1 bookmark": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks remove %v", PostIDExists)},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("Removed bookmark: [:link:](https://myhost.com//pl/ID2) `label1` `label2` Title2 - ")),
			expectedContains:  nil,
		},
		"REMOVE User successfully deletes 3 bookmark": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks remove %v,%v,%v", PostIDExists, b2ID, b3ID)},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: "",
			expectedContains: []string{
				"Removed bookmarks:",
				"[:link:](https://myhost.com//pl/ID2) `label1` `label2` Title2 - bookmarks initialized. Times created and same",
				"[:link:](https://myhost.com//pl/ID2) `label1` `label2` Title2 - bookmarks initialized. Times created and same",
				"[:link:](https://myhost.com//pl/ID3) Title3 - bookmarks already updated once",
			},
		},
	}
	for name, tt := range tests {
		api := makeAPIMock()
		tt.commandArgs.UserId = UserID
		siteURL := "https://myhost.com"
		api.On("GetPost", PostIDDoesNotExist).Return(nil, &model.AppError{Message: "An Error Occurred"})
		api.On("GetPost", b1ID).Return(&model.Post{Message: "this is the post.Message"}, nil)
		api.On("GetPost", b2ID).Return(&model.Post{Message: "this is the post.Message"}, nil)
		api.On("GetPost", b3ID).Return(&model.Post{Message: "this is the post.Message"}, nil)
		api.On("GetPost", b4ID).Return(&model.Post{Message: "this is the post.message"}, nil)
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
