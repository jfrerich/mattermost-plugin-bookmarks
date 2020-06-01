package main

// Bookmark contains information about an individual bookmark
type Bookmark struct {
	PostID     string   `json:"postid"`              // PostID is the ID for the bookmarked post and doubles as the Bookmark ID
	Title      string   `json:"title,omitempty"`     // Title given to the bookmark
	CreateAt   int64    `json:"create_at"`           // The original creation time of the bookmark
	ModifiedAt int64    `json:"update_at"`           // The original creation time of the bookmark
	LabelIDs   []string `json:"label_ids,omitempty"` // Array of labels added to the bookmark
}

func (bm *Bookmark) hasUserTitle() bool {
	return bm.getTitle() != ""
}

func (bm *Bookmark) hasLabels() bool {
	return bm.getLabelIDs() != nil
}

func (bm *Bookmark) getTitle() string {
	return bm.Title
}

func (bm *Bookmark) setTitle(title string) {
	bm.Title = title
}

func (bm *Bookmark) getLabelIDs() []string {
	return bm.LabelIDs
}

func (bm *Bookmark) addLabelIDs(ids []string) {
	bm.LabelIDs = ids
}
