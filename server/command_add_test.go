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

const addPrefixMsg = "Added bookmark: [:link:](https://myhost.com/_redirect/pl/"

func TestExecuteCommandAdd(t *testing.T) {
	tests := map[string]struct {
		commandArgs         *model.CommandArgs
		bookmarks           *Bookmarks
		labels              *Labels
		expectedMsgPrefix   string
		expectedContains    []string
		expectedNotContains []string
	}{
		"User doesn't provide an ID": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks add"},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace("Missing "),
			expectedContains:  []string{"Missing sub-command", "bookmarks add"},
		},
		"PostID doesn't exist": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v", PostIDDoesNotExist)},
			bookmarks:         nil,
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("PostID `%v` is not a valid postID", PostIDDoesNotExist)),
			expectedContains:  nil,
		},
		"Bookmark added  no title provided": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v", p1ID)},
			bookmarks:         getExecuteCommandTestBookmarks(),
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("%sID1) `TitleFromPost` this is the post.Message", addPrefixMsg)),
			expectedContains:  nil,
		},

		// TITLE PROVIDED; NO LABELS
		"Bookmark added  title provided no spaces": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v %v", PostIDExists, "TitleProvidedByUser")},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("%sID2)", addPrefixMsg)),
			expectedContains:  []string{"TitleProvidedByUser"},
		},
		"Bookmark added  title provided with spaces": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v %v", PostIDExists, "Title Provided By User")},
			bookmarks:         getExecuteCommandTestBookmarks(),
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("%sID2)", addPrefixMsg)),
			expectedContains:  []string{"Title Provided By User"},
		},

		// HAS TITLES AND LABELS
		"Bookmark added  title provided with spaces and labels": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v %v --labels %v", PostIDExists, "Title Provided By User", "label1,label2,label8")},
			bookmarks:         getExecuteCommandTestBookmarks(),
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("%sID2)", addPrefixMsg)),
			expectedContains:  []string{"label1", "label2", "label8", "Title Provided By User"},
		},
		"no flag optionBookmark added  title provided with spaces and labels": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v %v --labels %v", PostIDExists, "Title Provided By User", "label1,label2")},
			bookmarks:         getExecuteCommandTestBookmarks(),
			labels:            getExecuteCommandTestLabels(),
			expectedMsgPrefix: strings.TrimSpace(fmt.Sprintf("%sID2) `label1` `label2` Title Provided By User", addPrefixMsg)),
			expectedContains:  []string{"label1", "label2", "Title Provided By User"},
		},
		"Bookmark added  title provided with labels": {
			commandArgs:         &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v %v --labels label1,label2", PostIDExists, "TitleProvidedByUser")},
			bookmarks:           getExecuteCommandTestBookmarks(),
			labels:              getExecuteCommandTestLabels(),
			expectedMsgPrefix:   strings.TrimSpace(fmt.Sprintf("%sID2) ", addPrefixMsg)),
			expectedContains:    []string{"label1", "label2", "TitleProvidedByUser"},
			expectedNotContains: []string{"--labels"},
		},

		// HAS LABELS; NO TITLES
		"Bookmark unknown flag provided": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v --unknownflag", p1ID)},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace("Unable to parse options, unknown flag: --unknownflag"),
			expectedContains:  nil,
		},
		"Bookmark --labels provided without options": {
			commandArgs:       &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v --labels", p1ID)},
			bookmarks:         getExecuteCommandTestBookmarks(),
			expectedMsgPrefix: strings.TrimSpace("Unable to parse options, flag needs an argument: --labels"),
			expectedContains:  nil,
		},
		"Bookmark --labels provided with one label": {
			commandArgs:         &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v --labels label1", p1ID)},
			bookmarks:           getExecuteCommandTestBookmarks(),
			labels:              getExecuteCommandTestLabels(),
			expectedMsgPrefix:   strings.TrimSpace(fmt.Sprintf("%sID1", addPrefixMsg)),
			expectedNotContains: []string{"--labels"},
			expectedContains:    nil,
		},
		"Bookmark --labels provided with two labels": {
			commandArgs:         &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v --labels label1,label2", p1ID)},
			bookmarks:           getExecuteCommandTestBookmarks(),
			labels:              getExecuteCommandTestLabels(),
			expectedMsgPrefix:   strings.TrimSpace(fmt.Sprintf("%sID1", addPrefixMsg)),
			expectedContains:    []string{"label1", "label2"},
			expectedNotContains: []string{"--labels"},
		},
		"Bookmark add 3 labels two exist one new": {
			commandArgs:         &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v --labels label1,label2,label8", p1ID)},
			bookmarks:           getExecuteCommandTestBookmarks(),
			labels:              getExecuteCommandTestLabels(),
			expectedMsgPrefix:   strings.TrimSpace(fmt.Sprintf("%sID1", addPrefixMsg)),
			expectedContains:    []string{"label1", "label2", "label8"},
			expectedNotContains: []string{"--labels"},
		},
		"Bookmark added  test labels are sorted": {
			commandArgs: &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v --labels label1,l8,l2,aa,cc,bb,xx", p1ID)},
			// commandArgs: &model.CommandArgs{Command: fmt.Sprintf("/bookmarks add %v --labels label1,xx", p1ID)},
			bookmarks:           getExecuteCommandTestBookmarks(),
			labels:              getExecuteCommandTestLabels(),
			expectedMsgPrefix:   strings.TrimSpace(fmt.Sprintf("%sID1) `aa` `bb` `cc` `l2` `l8` `label1` `xx`", addPrefixMsg)),
			expectedContains:    []string{"l1", "l2", "l8"},
			expectedNotContains: []string{"--labels"},
			// expectedContains: nil,
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
		// api.On("ByID", mock.Anything).Return(true)

		jsonBmarks, err := json.Marshal(tt.bookmarks)
		api.On("KVGet", getBookmarksKey(tt.commandArgs.UserId)).Return(jsonBmarks, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)

		jsonLabels, err := json.Marshal(tt.labels)
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
				if tt.expectedNotContains != nil {
					for i := range tt.expectedNotContains {
						assert.NotContains(t, actual, tt.expectedNotContains[i])
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
