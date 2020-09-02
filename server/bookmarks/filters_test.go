package bookmarks

import (
	"sort"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi/mock_pluginapi"
	"github.com/stretchr/testify/assert"
)

func TestApplyFilters(t *testing.T) {
	// create some test bookmarks
	b1 := &Bookmark{
		PostID: "postID1",
		Title:  "This is my first title",
	}
	b2 := &Bookmark{
		PostID:   "postID2",
		LabelIDs: []string{"LID1", "LID2"},
		Title:    "This is my second title",
	}
	b3 := &Bookmark{
		PostID:   "postID3",
		LabelIDs: []string{"LID1", "LID2", "LID3"},
		Title:    "This is my third title",
	}

	// User1 has no bookmarks
	u1 := "userID1"
	bmarksU1 := NewBookmarks(u1)

	// User2 has 3 existing bookmarks
	u2 := "userID2"
	bmarksU2 := NewBookmarks(u2)
	bmarksU2.ByID[b1.PostID] = b1
	bmarksU2.ByID[b2.PostID] = b2
	bmarksU2.ByID[b3.PostID] = b3

	tests := []struct {
		name             string
		bmarks           *Bookmarks
		titleText        string
		labelIDs         []string
		expectedBmarkIDs []string
	}{
		{
			name:             "TITLE no bmarks  no title text found",
			bmarks:           bmarksU1,
			titleText:        "LIDNOTFOUND",
			expectedBmarkIDs: nil,
		},
		{
			name:             "TITLE no bmarks  no text requested",
			bmarks:           bmarksU1,
			titleText:        "",
			expectedBmarkIDs: nil,
		},
		{
			name:             "TITLE has bmarks  no title text found",
			titleText:        "LIDNOTFOUND",
			bmarks:           bmarksU1,
			expectedBmarkIDs: nil,
		},
		{
			name:             "TITLE has bmarks  no title text requested",
			titleText:        "",
			bmarks:           bmarksU2,
			expectedBmarkIDs: []string{"postID1", "postID2", "postID3"},
		},
		{
			name:             "TITLE has bmarks  title text requested  one found  1 ",
			titleText:        "first",
			bmarks:           bmarksU2,
			expectedBmarkIDs: []string{"postID1"},
		},
		{
			name:             "TITLE has bmarks  title text requested  one found  2",
			titleText:        "title",
			bmarks:           bmarksU2,
			expectedBmarkIDs: []string{"postID1", "postID2", "postID3"},
		},
		{
			name:             "LABELS no bmarks  no labels found",
			bmarks:           bmarksU1,
			labelIDs:         []string{"LIDNOTFOUND"},
			expectedBmarkIDs: nil,
		},
		{
			name:             "LABELS no bmarks  exist and no labels requested",
			bmarks:           bmarksU1,
			labelIDs:         nil,
			expectedBmarkIDs: nil,
		},
		{
			name:             "LABELS has bmarks  no labels found",
			labelIDs:         []string{"LIDNOTFOUND"},
			bmarks:           bmarksU2,
			expectedBmarkIDs: nil,
		},
		{
			name:             "LABELS has bmarks  no labels requested",
			labelIDs:         nil,
			bmarks:           bmarksU2,
			expectedBmarkIDs: []string{"postID1", "postID2", "postID3"},
		},
		{
			name:             "LABELS has bmarks  two labels requested",
			labelIDs:         []string{"LID1", "LID2"},
			bmarks:           bmarksU2,
			expectedBmarkIDs: []string{"postID2", "postID3"},
		},
		{
			name:             "LABELS has bmarks  one labels requested",
			labelIDs:         []string{"LID3"},
			bmarks:           bmarksU2,
			expectedBmarkIDs: []string{"postID3"},
		},
		{
			name:             "LABELS_TITLES has bmarks  no labels or titles requested",
			titleText:        "third title",
			labelIDs:         []string{"LID1", "LID2", "LID3"},
			bmarks:           bmarksU2,
			expectedBmarkIDs: []string{"postID3"},
		},
		{
			name:             "LABELS_TITLES has bmarks  one label match and one title match",
			titleText:        "third title",
			labelIDs:         []string{"LID3"},
			bmarks:           bmarksU2,
			expectedBmarkIDs: []string{"postID3"},
		},
		{
			name:             "LABELS_TITLES has bmarks  one label match but no title match",
			titleText:        "DOESNTMATCH",
			labelIDs:         []string{"LID3"},
			bmarks:           bmarksU2,
			expectedBmarkIDs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPluginAPI := mock_pluginapi.NewMockAPI(ctrl)

			bmarks := tt.bmarks
			bmarks.api = mockPluginAPI

			filters := &Filters{
				TitleText: tt.titleText,
				LabelIDs:  tt.labelIDs,
			}

			bmarks, err := bmarks.ApplyFilters(filters)

			assert.Nil(t, err)
			var ids []string
			for id := range bmarks.ByID {
				ids = append(ids, id)
			}
			sort.Strings(ids)
			sort.Strings(tt.expectedBmarkIDs)

			// check expected bmarks exist after filter
			assert.Equal(t, tt.expectedBmarkIDs, ids)
		})
	}
}
