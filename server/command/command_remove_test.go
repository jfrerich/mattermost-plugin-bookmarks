package command

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/bookmarks"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi/mock_pluginapi"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/assert"
)

func TestExecuteCommandRemove(t *testing.T) {
	tests := map[string]struct {
		command             string
		expectedMsgPrefix   string
		expectedContains    []string
		expectedNotContains []string
	}{
		"User does not provide an ID": {
			command:           "/bookmarks remove",
			expectedMsgPrefix: strings.TrimSpace("Missing "),
			expectedContains:  []string{"Missing sub-command", "bookmarks remove"},
		},
		"User tries to delete a bookmark but has none": {
			command:           "/bookmarks remove bmarkID",
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("Bookmark `bmarkID` does not exist")),
			expectedContains:  nil,
		},
		"User has bmarks tries to delete bookmark that doesnt exist": {
			command:           fmt.Sprintf("/bookmarks remove %v", PostIDDoesNotExist),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("Bookmark `%v` does not exist", PostIDDoesNotExist)),
			expectedContains:  nil,
		},
		"User successfully deletes 1 bookmark": {
			command:           fmt.Sprintf("/bookmarks remove %v", PostIDExists),
			expectedMsgPrefix: strings.TrimSpace("Removed bookmark: [:link:](https://myhost.com/_redirect/pl/ID2) `label1` `label2` **_Title2 - "),
			expectedContains:  nil,
		},
		"User successfully deletes 3 bookmark": {
			command:           fmt.Sprintf("/bookmarks remove %v %v %v", p1ID, p2ID, p3ID),
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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPluginAPI := mock_pluginapi.NewMockAPI(ctrl)

		mockPluginAPI.EXPECT().GetPost(p1ID).Return(&model.Post{Message: "this is the post.Message"}, nil).AnyTimes()
		mockPluginAPI.EXPECT().GetPost(p2ID).Return(&model.Post{Message: "this is the post.Message"}, nil).AnyTimes()
		mockPluginAPI.EXPECT().GetPost(p3ID).Return(&model.Post{Message: "this is the post.Message"}, nil).AnyTimes()

		bmarks := getExecuteCommandTestBookmarks()
		jsonBmarks, err := json.Marshal(bmarks)
		assert.Nil(t, err)

		labels := getExecuteCommandTestLabels()
		jsonLabels, err := json.Marshal(labels)
		assert.Nil(t, err)

		config := &model.Config{
			ServiceSettings: model.ServiceSettings{
				SiteURL: model.NewString("https://myhost.com"),
			},
		}
		mockPluginAPI.EXPECT().GetConfig().Return(config).AnyTimes()

		mockPluginAPI.EXPECT().KVSet(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		mockPluginAPI.EXPECT().KVGet(bookmarks.GetLabelsKey(UserID)).Return(jsonLabels, nil).AnyTimes()
		mockPluginAPI.EXPECT().KVGet(bookmarks.GetBookmarksKey(UserID)).Return(jsonBmarks, nil).AnyTimes()

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
