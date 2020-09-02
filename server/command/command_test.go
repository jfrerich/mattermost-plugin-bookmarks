package command

import (
	"encoding/json"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/bookmarks"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi/mock_pluginapi"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestPlugin struct {
	plugin.MattermostPlugin
}

//nolint
func makePlugin(api *plugintest.API) *TestPlugin {
	p := &TestPlugin{}
	p.SetAPI(api)
	return p
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
		command             string
		bookmarks           *bookmarks.Bookmarks
		expectedMsgPrefix   string
		expectedContains    []string
		expectedNotContains []string
	}{
		// No Slash Command
		"No slash command": {
			command:           "/bookmarks",
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("###### Bookmarks Slash Command Help "),
			expectedContains:  []string{"bookmarks add", "bookmarks view", "bookmarks remove"},
		},

		// Unknown Slash Command
		"UNKNOWN slash command": {
			command:           "/bookmarks UnknownCommand",
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("Unknown command: /bookmarks UnknownCommand"),
			expectedContains:  []string{},
		},

		// Help Slash Command
		"HELP slash command": {
			command:           "/bookmarks help",
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("###### Bookmarks Slash Command Help "),
			expectedContains:  []string{"bookmarks add", "bookmarks view", "bookmarks remove"},
		},
	}
	for name, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPluginAPI := mock_pluginapi.NewMockAPI(ctrl)

		api := makeAPIMock()
		// tt.command.UserId = UserID
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
		api.On("KVGet", bookmarks.GetBookmarksKey(UserID)).Return(jsonBmarks, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)

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

func makeAPIMock() *plugintest.API {
	api := &plugintest.API{}

	api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
	api.On("LogWarn", mock.Anything, mock.Anything, mock.Anything).Maybe()
	api.On("LogError", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

	return api
}
