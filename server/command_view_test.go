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

func getExecuteCommandViewBookmarks() *Bookmarks {
	api := makeAPIMock()
	api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
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
		CreateAt:   model.GetMillis() + 5,
		ModifiedAt: model.GetMillis(),
		LabelIDs:   []string{"UUID1", "UUID2", "UUID3"},
	}
	b3 := &Bookmark{
		PostID:     p3ID,
		Title:      b3Title,
		CreateAt:   model.GetMillis() + 2,
		ModifiedAt: model.GetMillis(),
		LabelIDs:   []string{"UUID3"},
	}
	b4 := &Bookmark{
		PostID:     p4ID,
		CreateAt:   model.GetMillis() + 3,
		ModifiedAt: model.GetMillis(),
	}

	_ = bmarks.add(b1)
	_ = bmarks.add(b2)
	_ = bmarks.add(b3)
	_ = bmarks.add(b4)

	return bmarks
}

func getExecuteCommandViewLabels() *Labels {
	l1 := &Label{Name: "label1"}
	l2 := &Label{Name: "label2"}
	l3 := &Label{Name: "label3"}

	api := makeAPIMock()
	api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
	labels := NewLabelsWithUser(api, UserID)
	_ = labels.add("UUID1", l1)
	_ = labels.add("UUID2", l2)
	_ = labels.add("UUID3", l3)

	return labels
}

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
		strings.TrimSpace(getLegendText()),
		"#### Bookmarks",
		"[:link:](https://myhost.com/_redirect/pl/ID1) `label1` `label2` **_Title1 - New Bookmark - times are zero_**",
		"[:link:](https://myhost.com/_redirect/pl/ID3) `label3` **_Title3 - bookmarks already updated once_**",
		"[:link:](https://myhost.com/_redirect/pl/ID4) **`TFP`** this is the post.Message",
		"[:link:](https://myhost.com/_redirect/pl/ID2) `label1` `label2` `label3` **_Title2 - bookmarks initialized. Times created and same_**",
	}

	tests := map[string]struct {
		commandArgs         *model.CommandArgs
		bookmarks           *Bookmarks
		expectedMsgPrefix   string
		expectedContains    []string
		expectedNotContains []string
	}{
		// User has no bookmarks
		"User has no bookmarks": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view"},
			bookmarks:         &Bookmarks{},
			expectedMsgPrefix: strings.TrimSpace("You do not have any saved bookmarks"),
			expectedContains:  nil,
		},
		"User has no bookmarks2": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view"},
			bookmarks:         &Bookmarks{},
			expectedMsgPrefix: strings.TrimSpace("You do not have any saved bookmarks"),
			expectedContains:  nil,
		},

		// View individual bookmark
		"User requests to view bookmark by ID that has a title defined": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view ID2"},
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
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view"},
			expectedMsgPrefix: strings.TrimSpace(getLegendText()),
			expectedContains:  []string{"Bookmarks", "ID1", "ID2", "ID3"},
		},
		"User has 4 bookmarks  All with titles  One without": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view"},
			expectedMsgPrefix: strings.TrimSpace(getLegendText()),
			expectedContains:  defaultSortString,
		},

		// View Sorting
		"Sorted by createdAt - default sort": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view"},
			expectedMsgPrefix: strings.TrimSpace(strings.Join(defaultSortString, "\n")),
			expectedContains:  nil,
		},

		// filter bookmarks
		"User filter by label  filter one label  label1": {
			commandArgs:       &model.CommandArgs{Command: "/bookmarks view --filter-labels label1"},
			expectedMsgPrefix: strings.TrimSpace(getLegendText()),
			expectedContains:  []string{"Bookmarks", "ID1", "ID2"},
		},
		"User filter by label  filter two labels": {
			commandArgs:         &model.CommandArgs{Command: "/bookmarks view --filter-labels label1,label2"},
			expectedMsgPrefix:   strings.TrimSpace(getLegendText()),
			expectedContains:    []string{"Bookmarks", "ID1", "ID2"},
			expectedNotContains: []string{"ID3", "ID4"},
		},
		"User filter by label  filter one label  label3": {
			commandArgs:         &model.CommandArgs{Command: "/bookmarks view --filter-labels label3"},
			expectedMsgPrefix:   strings.TrimSpace(getLegendText()),
			expectedContains:    []string{"Bookmarks", "ID3", "ID2"},
			expectedNotContains: []string{"ID1", "ID4"},
		},
		"User filter by label  filter all available labels": {
			commandArgs:         &model.CommandArgs{Command: "/bookmarks view --filter-labels label1,label2,label3"},
			expectedMsgPrefix:   strings.TrimSpace(getLegendText()),
			expectedContains:    []string{"Bookmarks", "ID1", "ID2", "ID3"},
			expectedNotContains: []string{"ID4"},
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

		bookmarks := getExecuteCommandViewBookmarks()
		if tt.bookmarks != nil {
			bookmarks = tt.bookmarks
		}
		jsonBmarks, err := json.Marshal(bookmarks)

		// jsonBmarks, err = json.Marshal(tt.bookmarks)
		api.On("KVGet", getBookmarksKey(tt.commandArgs.UserId)).Return(jsonBmarks, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)

		labels := getExecuteCommandViewLabels()
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
