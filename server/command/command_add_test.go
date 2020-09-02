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

const addPrefixMsg = "Added bookmark: [:link:](https://myhost.com/_redirect/pl/"

func TestExecuteCommandAdd(t *testing.T) {
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
		command             string
		bookmarks           *bookmarks.Bookmarks
		labels              *bookmarks.Labels
		expectedMsgPrefix   string
		expectedContains    []string
		expectedNotContains []string
	}{
		"User doesn't provide an ID": {
			command:           "/bookmarks add",
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("Missing "),
			expectedContains:  []string{"Missing sub-command", "bookmarks add"},
		},
		"PostID doesn't exist": {
			command:           fmt.Sprintf("/bookmarks add %v", PostIDDoesNotExist),
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("PostID `%v` is not a valid postID", PostIDDoesNotExist)),
			expectedContains:  nil,
		},
		"Bookmark added  no title provided": {
			command:           fmt.Sprintf("/bookmarks add %v", p1ID),
			bookmarks:         getExecuteCommandTestBookmarks(),
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("%sID1) **`TFP`** this is the post.Message", addPrefixMsg)),
			expectedContains:  nil,
		},

		// TITLE PROVIDED; NO LABELS
		"Bookmark added  title provided no spaces": {
			command:           fmt.Sprintf("/bookmarks add %v %v", PostIDExists, "TitleProvidedByUser"),
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("%sID2)", addPrefixMsg)),
			expectedContains:  []string{"TitleProvidedByUser"},
		},
		"Bookmark added  title provided with spaces": {
			command:           fmt.Sprintf("/bookmarks add %v %v", PostIDExists, "Title Provided By User"),
			bookmarks:         getExecuteCommandTestBookmarks(),
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("%sID2)", addPrefixMsg)),
			expectedContains:  []string{"Title Provided By User"},
		},

		// HAS TITLES AND LABELS
		"Bookmark added  title provided with spaces and labels": {
			command:           fmt.Sprintf("/bookmarks add %v %v --labels %v", PostIDExists, "Title Provided By User", "label1,label2,label8"),
			bookmarks:         getExecuteCommandTestBookmarks(),
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("%sID2)", addPrefixMsg)),
			expectedContains:  []string{"label1", "label2", "label8", "Title Provided By User"},
		},
		"no flag optionBookmark added  title provided with spaces and labels": {
			command:           fmt.Sprintf("/bookmarks add %v %v --labels %v", PostIDExists, "Title Provided By User", "label1,label2"),
			bookmarks:         getExecuteCommandTestBookmarks(),
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("%sID2) `label1` `label2` **_Title Provided By User_**", addPrefixMsg)),
			expectedContains:  []string{"label1", "label2", "Title Provided By User"},
		},
		"Bookmark added  title provided with labels": {
			command:             fmt.Sprintf("/bookmarks add %v %v --labels label1,label2", PostIDExists, "TitleProvidedByUser"),
			bookmarks:           getExecuteCommandTestBookmarks(),
			labels:              getExecuteCommandTestLabels(),
			expectedMsgPrefix:   strings.TrimSpace(fmt.Sprintf("%sID2) ", addPrefixMsg)),
			expectedContains:    []string{"label1", "label2", "TitleProvidedByUser"},
			expectedNotContains: []string{"--labels"},
		},

		// HAS LABELS; NO TITLES
		"Bookmark unknown flag provided": {
			command:           fmt.Sprintf("/bookmarks add %v --unknownflag", p1ID),
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace("Unable to parse options, unknown flag: --unknownflag"),
			expectedContains:  nil,
		},
		"Bookmark --labels provided without options": {
			command:           fmt.Sprintf("/bookmarks add %v --labels", p1ID),
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace("Unable to parse options, flag needs an argument: --labels"),
			expectedContains:  nil,
		},
		"Bookmark --labels provided with one label": {
			command:             fmt.Sprintf("/bookmarks add %v --labels label1", p1ID),
			bookmarks:           getExecuteCommandTestBookmarks(),
			labels:              getExecuteCommandTestLabels(),
			expectedMsgPrefix:   strings.TrimSpace(fmt.Sprintf("%sID1", addPrefixMsg)),
			expectedNotContains: []string{"--labels"},
			expectedContains:    nil,
		},
		"Bookmark --labels provided with two labels": {
			command:             fmt.Sprintf("/bookmarks add %v --labels label1,label2", p1ID),
			bookmarks:           getExecuteCommandTestBookmarks(),
			labels:              getExecuteCommandTestLabels(),
			expectedMsgPrefix:   strings.TrimSpace(fmt.Sprintf("%sID1", addPrefixMsg)),
			expectedContains:    []string{"label1", "label2"},
			expectedNotContains: []string{"--labels"},
		},
		"Bookmark add 3 labels two exist one new": {
			command:             fmt.Sprintf("/bookmarks add %v --labels label1,label2,label8", p1ID),
			bookmarks:           getExecuteCommandTestBookmarks(),
			labels:              getExecuteCommandTestLabels(),
			expectedMsgPrefix:   strings.TrimSpace(fmt.Sprintf("%sID1", addPrefixMsg)),
			expectedContains:    []string{"label1", "label2", "label8"},
			expectedNotContains: []string{"--labels"},
		},
		"Bookmark added  test labels are sorted": {
			command:             fmt.Sprintf("/bookmarks add %v --labels label1,l8,l2,aa,cc,bb,xx", p1ID),
			bookmarks:           getExecuteCommandTestBookmarks(),
			labels:              getExecuteCommandTestLabels(),
			expectedMsgPrefix:   strings.TrimSpace(fmt.Sprintf("%sID1) **`TFP`** `aa` `bb` `cc` `l2` `l8` `label1` `xx`", addPrefixMsg)),
			expectedContains:    []string{"l1", "l2", "l8"},
			expectedNotContains: []string{"--labels"},
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

		mockPluginAPI.EXPECT().GetConfig().Return(config).AnyTimes()
		mockPluginAPI.EXPECT().GetPost(PostIDDoesNotExist).Return(nil, &model.AppError{Message: "An Error Occurred"}).AnyTimes()
		mockPluginAPI.EXPECT().GetPost(p1ID).Return(p1IDmodel, nil).AnyTimes()
		mockPluginAPI.EXPECT().GetPost(p2ID).Return(p2IDmodel, nil).AnyTimes()
		mockPluginAPI.EXPECT().GetPost(p3ID).Return(p3IDmodel, nil).AnyTimes()
		mockPluginAPI.EXPECT().GetPost(p4ID).Return(p4IDmodel, nil).AnyTimes()

		jsonBmarks, err := json.Marshal(tt.bookmarks)
		jsonLabels, err := json.Marshal(tt.labels)

		mockPluginAPI.EXPECT().KVGet(bookmarks.GetLabelsKey(UserID)).Return(jsonLabels, nil).AnyTimes()
		mockPluginAPI.EXPECT().KVGet(bookmarks.GetBookmarksKey(UserID)).Return(jsonBmarks, nil).AnyTimes()

		mockPluginAPI.EXPECT().KVSet(bookmarks.GetBookmarksKey(UserID), gomock.Any()).Return(nil).AnyTimes()
		mockPluginAPI.EXPECT().KVSet(bookmarks.GetLabelsKey(UserID), gomock.Any()).Return(nil).AnyTimes()

		t.Run(name, func(t *testing.T) {
			assert.Nil(t, err)

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
