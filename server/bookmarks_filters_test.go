package main

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApplyFilters(t *testing.T) {
	api := makeAPIMock()
	api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
	p := makePlugin(api)

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
	bmarksU1 := NewBookmarksWithUser(p.API, u1)

	// User2 has 3 existing bookmarks
	u2 := "userID2"
	bmarksU2 := NewBookmarksWithUser(p.API, u2)
	err := bmarksU2.add(b1)
	assert.Nil(t, err)
	err = bmarksU2.add(b2)
	assert.Nil(t, err)
	err = bmarksU2.add(b3)
	assert.Nil(t, err)

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
			expectedBmarkIDs: nil,
		},
		{
			name:             "TITLE has bmarks  no title text requested",
			titleText:        "",
			expectedBmarkIDs: []string{"postID1", "postID2", "postID3"},
		},
		{
			name:             "TITLE has bmarks  title text requested  one found",
			titleText:        "first",
			expectedBmarkIDs: []string{"postID1"},
		},
		{
			name:             "TITLE has bmarks  title text requested  one found",
			titleText:        "title",
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
			expectedBmarkIDs: nil,
		},
		{
			name:             "LABELS has bmarks  no labels requested",
			labelIDs:         nil,
			expectedBmarkIDs: []string{"postID1", "postID2", "postID3"},
		},
		{
			name:             "LABELS has bmarks  two labels requested",
			labelIDs:         []string{"LID1", "LID2"},
			expectedBmarkIDs: []string{"postID2", "postID3"},
		},
		{
			name:             "LABELS has bmarks  one labels requested",
			labelIDs:         []string{"LID3"},
			expectedBmarkIDs: []string{"postID3"},
		},
		{
			name:             "LABELS_TITLES has bmarks  no labels or titles requested",
			titleText:        "third title",
			labelIDs:         []string{"LID1", "LID2", "LID3"},
			expectedBmarkIDs: []string{"postID3"},
		},
		{
			name:             "LABELS_TITLES has bmarks  one label match and one title match",
			titleText:        "third title",
			labelIDs:         []string{"LID3"},
			expectedBmarkIDs: []string{"postID3"},
		},
		{
			name:             "LABELS_TITLES has bmarks  one label match but no title match",
			titleText:        "DOESNTMATCH",
			labelIDs:         []string{"LID3"},
			expectedBmarkIDs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bmarks := bmarksU2
			if tt.bmarks != nil {
				bmarks = tt.bmarks
			}

			filters := &BookmarksFilters{
				TitleText: tt.titleText,
				LabelIDs:  tt.labelIDs,
			}

			bmarks, err = bmarks.applyFilters(filters)
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
