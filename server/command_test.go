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

const (
	PostIDDoesNotExist = "PostIDDoesNotExist"
	PostIDExists       = "ID2"
	UserID             = "UserID"
	teamID1            = "teamID1"

	p1ID = "ID1"
	p2ID = "ID2"
	p3ID = "ID3"
	p4ID = "ID4"

	b1Title = "Title1 - New Bookmark - times are zero"
	b2Title = "Title2 - bookmarks initialized. Times created and same"
	b3Title = "Title3 - bookmarks already updated once"
)

func getExecuteCommandTestBookmarks() *Bookmarks {
	api := makeAPIMock()
	p := makePlugin(api)
	bmarks := NewBookmarksWithUser(p.API, UserID)

	b1 := &Bookmark{
		PostID:   p1ID,
		Title:    b1Title,
		LabelIDs: []string{"UUID1", "UUID2"},
	}
	b2 := &Bookmark{
		PostID:     p2ID,
		Title:      b2Title,
		CreateAt:   model.GetMillis(),
		ModifiedAt: model.GetMillis(),
		LabelIDs:   []string{"UUID1", "UUID2"},
	}
	b3 := &Bookmark{
		PostID:     p3ID,
		Title:      b3Title,
		CreateAt:   model.GetMillis(),
		ModifiedAt: model.GetMillis(),
	}
	b4 := &Bookmark{
		PostID:     p4ID,
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

	labels := NewLabels(api)
	labels.add("UUID1", l1)

	return bmarks
}

func TestExecuteCommand(t *testing.T) {
	p1IDmodel := &model.Post{
		Message:  "this is the post.Message",
		CreateAt: model.GetMillis(),
	}
	p2IDmodel := &model.Post{
		Message:  "this is the post.Message",
		CreateAt: model.GetMillis() + 1,
	}
	p3IDmodel := &model.Post{
		Message:  "this is the post.Message",
		CreateAt: model.GetMillis() + 2,
	}
	p4IDmodel := &model.Post{
		Message:  "this is the post.Message",
		CreateAt: model.GetMillis() + 3,
	}

	tests := map[string]struct {
		commandArgs       *model.CommandArgs
		bookmarks         *Bookmarks
		expectedMsgPrefix string
		expectedContains  []string
	}{
		// No Slash Command
		"No slash command": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("###### Bookmarks Slash Command Help "),
			expectedContains:  []string{"bookmarks add", "bookmarks view", "bookmarks remove"},
		},

		// Unknown Slash Command
		"UNKNOWN slash command": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks UnknownCommand"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("Unknown command: /bookmarks UnknownCommand"),
			expectedContains:  []string{},
		},

		// Help Slash Command
		"HELP slash command": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks help"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("###### Bookmarks Slash Command Help "),
			expectedContains:  []string{"bookmarks add", "bookmarks view", "bookmarks remove"},
		},
	}
	for name, tt := range tests {
		api := makeAPIMock()
		tt.commandArgs.UserId = UserID
		siteURL := "https://myhost.com"
		api.On("GetPost", PostIDDoesNotExist).Return(nil, &model.AppError{Message: "An Error Occurred"})
		api.On("GetPost", p1ID).Return(p1IDmodel, nil)
		api.On("GetPost", p2ID).Return(p2IDmodel, nil)
		api.On("GetPost", p3ID).Return(p3IDmodel, nil)
		api.On("GetPost", p4ID).Return(p4IDmodel, nil)
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
