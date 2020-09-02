package command

import (
	"encoding/json"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/bookmarks"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi/mock_pluginapi"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/utils"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/assert"
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

	defaultSortString := []string{
		strings.TrimSpace(utils.GetLegendText()),
		"#### Bookmarks",
		"[:link:](https://myhost.com/_redirect/pl/ID1) `label1` `label2` **_Title1 - New Bookmark - times are zero_**",
		"[:link:](https://myhost.com/_redirect/pl/ID3) `label3` **_Title3 - bookmarks already updated once_**",
		"[:link:](https://myhost.com/_redirect/pl/ID4) **`TFP`** this is the post.Message",
		"[:link:](https://myhost.com/_redirect/pl/ID2) `label1` `label2` `label3` **_Title2 - bookmarks initialized. Times created and same_**",
	}

	tests := map[string]struct {
		command             string
		bmarks              *bookmarks.Bookmarks
		expectedMsgPrefix   string
		expectedContains    []string
		expectedNotContains []string
	}{
		// User has no bookmarks
		"User has no bookmarks": {
			command:           "/bookmarks view",
			bmarks:            &bookmarks.Bookmarks{},
			expectedMsgPrefix: strings.TrimSpace("You do not have any saved bookmarks"),
			expectedContains:  nil,
		},

		// View individual bookmark
		"User requests to view bookmark by ID that has a title defined": {
			command:           "/bookmarks view ID2",
			expectedMsgPrefix: "",
			expectedContains: []string{
				"#### Bookmark Title [:link:](https://myhost.com/_redirect/pl/ID2)",
				"`label1` `label2`",
				"**Title2 - bookmarks initialized. Times created and same**",
				"##### Post Message",
				"this is the post.Message",
			},
		},

		// View all bookmarks
		"User has 3 bookmarks  All with titles provided": {
			command:           "/bookmarks view",
			expectedMsgPrefix: strings.TrimSpace(utils.GetLegendText()),
			expectedContains:  []string{"Bookmarks", "ID1", "ID2", "ID3"},
		},
		"User has 4 bookmarks  All with titles  One without": {
			command:           "/bookmarks view",
			expectedMsgPrefix: strings.TrimSpace(utils.GetLegendText()),
			expectedContains:  defaultSortString,
		},

		// View Sorting
		"Sorted by createdAt - default sort": {
			command:           "/bookmarks view",
			expectedMsgPrefix: strings.TrimSpace(strings.Join(defaultSortString, "\n")),
			expectedContains:  nil,
		},

		// filter bookmarks
		"User filter by label  filter one label  label1": {
			command:           "/bookmarks view --filter-labels label1",
			expectedMsgPrefix: strings.TrimSpace(utils.GetLegendText()),
			expectedContains:  []string{"Bookmarks", "ID1", "ID2"},
		},
		"User filter by label  filter two labels": {
			command:             "/bookmarks view --filter-labels label1,label2",
			expectedMsgPrefix:   strings.TrimSpace(utils.GetLegendText()),
			expectedContains:    []string{"Bookmarks", "ID1", "ID2"},
			expectedNotContains: []string{"ID3", "ID4"},
		},
		"User filter by label  filter one label  label3": {
			command:             "/bookmarks view --filter-labels label3",
			expectedMsgPrefix:   strings.TrimSpace(utils.GetLegendText()),
			expectedContains:    []string{"Bookmarks", "ID3", "ID2"},
			expectedNotContains: []string{"ID1", "ID4"},
		},
		"User filter by label  filter all available labels": {
			command:             "/bookmarks view --filter-labels label1,label2,label3",
			expectedMsgPrefix:   strings.TrimSpace(utils.GetLegendText()),
			expectedContains:    []string{"Bookmarks", "ID1", "ID2", "ID3"},
			expectedNotContains: []string{"ID4"},
		},
	}
	for name, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPluginAPI := mock_pluginapi.NewMockAPI(ctrl)

		config := &model.Config{
			ServiceSettings: model.ServiceSettings{
				SiteURL: model.NewString("https://myhost.com"),
			},
		}
		mockPluginAPI.EXPECT().GetConfig().Return(config).AnyTimes()
		mockPluginAPI.EXPECT().GetPost(PostIDDoesNotExist).Return(nil, &model.AppError{Message: "An Error Occurred"}).AnyTimes()
		mockPluginAPI.EXPECT().GetPost(p1ID).Return(p1IDmodel, nil).AnyTimes()
		mockPluginAPI.EXPECT().GetPost(p2ID).Return(p2IDmodel, nil).AnyTimes()
		mockPluginAPI.EXPECT().GetPost(p3ID).Return(p3IDmodel, nil).AnyTimes()
		mockPluginAPI.EXPECT().GetPost(p4ID).Return(p4IDmodel, nil).AnyTimes()

		bmarks := tt.bmarks
		if tt.bmarks == nil {
			bmarks = getExecuteCommandViewBookmarks()
		}

		jsonBmarks, err := json.Marshal(bmarks)
		assert.Nil(t, err)

		labels := getExecuteCommandViewLabels()
		jsonLabels, err := json.Marshal(labels)
		assert.Nil(t, err)

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
