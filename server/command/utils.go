package command

import (
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/bookmarks"
	"github.com/mattermost/mattermost-server/model"
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

func getExecuteCommandTestBookmarks() *bookmarks.Bookmarks {
	b1 := &bookmarks.Bookmark{
		PostID:   p1ID,
		Title:    b1Title,
		LabelIDs: []string{"UUID1", "UUID2"},
	}
	b2 := &bookmarks.Bookmark{
		PostID:     p2ID,
		Title:      b2Title,
		CreateAt:   model.GetMillis() + 5,
		ModifiedAt: model.GetMillis(),
		LabelIDs:   []string{"UUID1", "UUID2"},
	}
	b3 := &bookmarks.Bookmark{
		PostID:     p3ID,
		Title:      b3Title,
		CreateAt:   model.GetMillis() + 2,
		ModifiedAt: model.GetMillis(),
	}
	b4 := &bookmarks.Bookmark{
		PostID:     p4ID,
		CreateAt:   model.GetMillis() + 3,
		ModifiedAt: model.GetMillis(),
	}

	bmarks := bookmarks.NewBookmarks(UserID)
	bmarks.ByID[b1.PostID] = b1
	bmarks.ByID[b2.PostID] = b2
	bmarks.ByID[b3.PostID] = b3
	bmarks.ByID[b4.PostID] = b4

	return bmarks
}

func getExecuteCommandTestLabels() *bookmarks.Labels {
	l1 := &bookmarks.Label{
		Name: "label1",
	}
	l2 := &bookmarks.Label{
		Name: "label2",
	}
	l3 := &bookmarks.Label{
		Name: "label8",
	}

	labels := bookmarks.NewLabels(UserID)
	labels.ByID["UUID1"] = l1
	labels.ByID["UUID2"] = l2
	labels.ByID["UUID3"] = l3
	return labels
}

func getExecuteCommandViewBookmarks() *bookmarks.Bookmarks {
	b1 := &bookmarks.Bookmark{
		PostID:   p1ID,
		Title:    b1Title,
		LabelIDs: []string{"UUID1", "UUID2"},
	}
	b2 := &bookmarks.Bookmark{
		PostID:     p2ID,
		Title:      b2Title,
		CreateAt:   model.GetMillis() + 5,
		ModifiedAt: model.GetMillis(),
		LabelIDs:   []string{"UUID1", "UUID2", "UUID3"},
	}
	b3 := &bookmarks.Bookmark{
		PostID:     p3ID,
		Title:      b3Title,
		CreateAt:   model.GetMillis() + 2,
		ModifiedAt: model.GetMillis(),
		LabelIDs:   []string{"UUID3"},
	}
	b4 := &bookmarks.Bookmark{
		PostID:     p4ID,
		CreateAt:   model.GetMillis() + 3,
		ModifiedAt: model.GetMillis(),
	}

	bmarks := bookmarks.NewBookmarks(UserID)
	bmarks.ByID[b1.PostID] = b1
	bmarks.ByID[b2.PostID] = b2
	bmarks.ByID[b3.PostID] = b3
	bmarks.ByID[b4.PostID] = b4
	return bmarks
}

func getExecuteCommandViewLabels() *bookmarks.Labels {
	l1 := &bookmarks.Label{Name: "label1", ID: "UUID1"}
	l2 := &bookmarks.Label{Name: "label2", ID: "UUID2"}
	l3 := &bookmarks.Label{Name: "label3", ID: "UUID3"}

	labels := bookmarks.NewLabels(UserID)
	labels.ByID["UUID1"] = l1
	labels.ByID["UUID2"] = l2
	labels.ByID["UUID3"] = l3

	return labels
}
