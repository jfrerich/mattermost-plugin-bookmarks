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

func TestExecuteCommandView(t *testing.T) {
	p1IDmodel := &model.Post{
		Message:  "this is the post.Message",
		CreateAt: model.GetMillis(),
	}
	p2IDmodel := &model.Post{
		Message:  "this is the post.Message",
		CreateAt: model.GetMillis() + 5,
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
		// USER HAS NO BOOKMARKS
		"User has no bookmarks": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("You do not have any saved bookmarks"),
			expectedContains:  nil,
		},
		"User has no bookmarks2": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("You do not have any saved bookmarks"),
			expectedContains:  nil,
		},

		// VIEW INDIVIDUAL BOOKMARK
		"User requests to view bookmark by ID that has a title defined": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view ID2"},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: "",
			expectedContains: []string{
				"#### Bookmark Title [:link:](https://myhost.com//pl/ID2)",
				"`label1` `label2`",
				"**Title2 - bookmarks initialized. Times created and same**",
				"##### Post Message",
				"this is the post.Message",
			},
		},

		// VIEW ALL BOOKMARKS
		"User has 3 bookmarks  All with titles provided": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view"},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace("#### Bookmarks List"),
			expectedContains:  []string{"Bookmarks List", "ID1", "ID2", "ID3"},
		},
		"User has 4 bookmarks  All with titles  One without": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view"},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace("#### Bookmarks List"),
			expectedContains: []string{
				"Bookmarks List",
				"[:link:](https://myhost.com//pl/ID1) `label1` `label2` Title1 - New Bookmark - times are zero",
				"[:link:](https://myhost.com//pl/ID2) `label1` `label2` Title2 - bookmarks initialized. Times created and same",
				"[:link:](https://myhost.com//pl/ID3) Title3 - bookmarks already updated once",
				"[:link:](https://myhost.com//pl/ID4) `TitleFromPost`"},
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

		labels := getExecuteCommandTestLabels()
		jsonLabels, err := json.Marshal(labels)
		api.On("KVGet", getLabelsKey(tt.commandArgs.UserId)).Return(jsonLabels, nil)

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
