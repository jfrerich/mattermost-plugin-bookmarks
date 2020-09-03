package bookmarks

import (
	"fmt"
	"sort"

	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/utils"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

type IBookmarks interface {
	GetBookmark(bmarkID string) (*Bookmark, error)
	AddBookmark(bmark *Bookmark) error
	DeleteBookmark(bmarkID string) error
	DeleteLabel(bmarkID string, labelID string) error
	GetBookmarksWithLabelID(labelID string) (*Bookmarks, error)
	GetBmarkTextOneLine(bmark *Bookmark, labelNames []string) (string, error)
	ApplyFilters(filters *Filters) (*Bookmarks, error)
	ByPostCreateAt() ([]*Bookmark, error)
	GetBmarkLabelNames(bmark *Bookmark) ([]string, error)
}

// Bookmarks contains a map of bookmarks
type Bookmarks struct {
	ByID   map[string]*Bookmark
	api    pluginapi.API
	userID string
}

func NewBookmarks(userID string) *Bookmarks {
	return &Bookmarks{
		ByID:   make(map[string]*Bookmark),
		userID: userID,
	}
}

// NewBookmarksWithUser returns an initialized Bookmarks for a User
func NewBookmarksWithUser(api pluginapi.API, userID string) (*Bookmarks, error) {
	bb, appErr := api.KVGet(GetBookmarksKey(userID))
	if appErr != nil {
		return nil, errors.Wrapf(appErr, "Unable to get bookmarks for user %s", userID)
	}

	userBmarks, err := BookmarksFromJSON(bb)
	if err != nil {
		return nil, err
	}
	userBmarks.api = api
	userBmarks.userID = userID

	return userBmarks, nil
}

func (b *Bookmarks) exists(bmarkID string) (*Bookmark, bool) {
	if bmark, ok := b.ByID[bmarkID]; ok {
		return bmark, true
	}
	return nil, false
}

// getBookmark returns a bookmark with the specified bookmarkID
func (b *Bookmarks) GetBookmark(bmarkID string) (*Bookmark, error) {
	if b == nil {
		return nil, nil
	}
	bmark, ok := b.exists(bmarkID)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Bookmark `%v` does not exist", bmarkID))
	}
	return bmark, nil
}

func (b *Bookmarks) updateTimes(bmarkID string) *Bookmark {
	bmark, _ := b.GetBookmark(bmarkID)
	if bmark.CreateAt == 0 {
		bmark.CreateAt = model.GetMillis()
		bmark.ModifiedAt = bmark.CreateAt
	}
	bmark.ModifiedAt = model.GetMillis()
	return bmark
}

// addBookmark stores the bookmark in a map,
func (b *Bookmarks) AddBookmark(bmark *Bookmark) error {
	// bookmark already exists, update ModifiedAt and save
	_, ok := b.exists(bmark.PostID)
	if ok {
		b.updateTimes(bmark.PostID)
		b.updateLabels(bmark)
	}

	// Add or update the bookmark
	b.ByID[bmark.PostID] = bmark
	if err := b.StoreBookmarks(); err != nil {
		return errors.Wrap(err, "failed to add bookmark")
	}
	return nil
}

// ByPostCreateAt returns an array of bookmarks sorted by post.CreateAt times
func (b *Bookmarks) ByPostCreateAt() ([]*Bookmark, error) {
	// build temp map
	tempMap := make(map[int64]string)
	for _, bmark := range b.ByID {
		post, appErr := b.api.GetPost(bmark.PostID)
		if appErr != nil {
			return nil, appErr
		}
		tempMap[post.CreateAt] = bmark.PostID
	}

	// sort post.CreateAt (keys)
	keys := make([]int, 0, len(tempMap))
	for k := range tempMap {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	// reconstruct the bookmarks in a sorted array
	var bookmarks []*Bookmark
	for _, k := range keys {
		bmark := b.ByID[tempMap[int64(k)]]
		bookmarks = append(bookmarks, bmark)
	}

	return bookmarks, nil
}

// func (b *Bookmarks) GetBookmarksWithLabelID(labelID string) (IBookmarks, error) {
func (b *Bookmarks) GetBookmarksWithLabelID(id string) (*Bookmarks, error) {
	// FIXME: This should not require setting the api again.
	bmarks := NewBookmarks(b.userID)
	bmarks.api = b.api

	for _, bmark := range b.ByID {
		if bmark.hasLabels() {
			for _, lid := range bmark.GetLabelIDs() {
				if id == lid {
					if err := bmarks.AddBookmark(bmark); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return bmarks, nil
}

// DeleteLabel deletes a label from a bookmark
func (b *Bookmarks) DeleteLabel(bmarkID string, labelID string) error {
	bmark, err := b.GetBookmark(bmarkID)
	if err != nil {
		return err
	}

	var newLabels []string
	origLabels := bmark.GetLabelIDs()
	for _, ID := range origLabels {
		if labelID == ID {
			continue
		}
		newLabels = append(newLabels, ID)
	}

	bmark.AddLabelIDs(newLabels)
	if err := b.AddBookmark(bmark); err != nil {
		return err
	}

	return nil
}

func (b *Bookmarks) updateLabels(bmark *Bookmark) *Bookmark {
	bmarkOrig, _ := b.GetBookmark(bmark.PostID)
	bmarkOrig.AddLabelIDs(bmark.GetLabelIDs())
	return bmark
}

// GetBmarkLabelNames returns an array of labelNames for a given bookmark
func (b *Bookmarks) GetBmarkLabelNames(bmark *Bookmark) ([]string, error) {
	labels, err := NewLabelsWithUser(b.api, b.userID)
	if err != nil {
		return nil, err
	}

	var labelNames []string
	for _, id := range bmark.GetLabelIDs() {
		name, err := labels.GetNameFromID(id)
		if err != nil {
			return nil, err
		}
		labelNames = append(labelNames, name)
	}
	return labelNames, nil
}

// GetBmarkTextOneLine returns a single line bookmark text used for an ephemeral post
func (b *Bookmarks) GetBmarkTextOneLine(bmark *Bookmark, labelNames []string) (string, error) {
	postMessage, err := b.getTitleFromPost(bmark.PostID)
	if err != nil {
		return "", err
	}

	codeBlockedNames := GetCodeBlockedLabels(labelNames)

	// bold and italicize titles saved by the user
	title := "**_" + bmark.GetTitle() + "_**"

	if !bmark.HasUserTitle() {
		// display the first portion of the post message in place of a title
		title = postMessage
		// prepend the title from post label before other labels
		codeBlockedNames = " " + utils.TitleFromPostLabel + codeBlockedNames
	}

	text := fmt.Sprintf("%s%s %s\n", getIconLink(b.api, bmark.PostID), codeBlockedNames, title)

	return text, nil
}

// getTitleFromPost returns a title generated from a Post.Message
func (b *Bookmarks) getTitleFromPost(postID string) (string, error) {
	// MaxTitleCharacters is the maximum length of characters displayed in a
	// bookmark title
	// MaxTitleCharacters = 30

	// TODO: set limit to number of character from post.Message
	// numChars := math.Min(float64(len(post.Message)), MaxTitleCharacters)
	// bookmark.Title = post.Message[0:int(numChars)]

	post, appErr := b.api.GetPost(postID)
	if appErr != nil {
		return "", appErr
	}
	title := post.Message
	return title, nil
}

// GetCodeBlockedLabels returns a list of individually codeblocked names
func GetCodeBlockedLabels(names []string) string {
	labels := ""
	sort.Strings(names)
	for _, name := range names {
		labels += fmt.Sprintf(" `%s`", name)
	}
	return labels
}

// getBmarksEphemeralText returns a the text for posting all bookmarks in an
// ephemeral message
func (b *Bookmarks) GetBmarksEphemeralText(userID string, filters *Filters) (string, error) {
	var err error
	if filters != nil {
		b, err = b.ApplyFilters(filters)
		if err != nil {
			return "", err
		}
	}

	// bookmarks is nil if user has never added a bookmark.
	// bookmarks.ByID will be empty if user created a bookmark and then deleted
	// it and now has 0 bookmarks
	if b == nil || len(b.ByID) == 0 {
		return "You do not have any saved bookmarks", nil
	}

	bmarksSorted, err := b.ByPostCreateAt()
	if err != nil {
		return "", err
	}

	text := utils.GetLegendText()
	text += "#### Bookmarks\n"
	for _, bmark := range bmarksSorted {
		labelNames, err := b.GetBmarkLabelNames(bmark)
		if err != nil {
			return "", err
		}
		nextText, err := b.GetBmarkTextOneLine(bmark, labelNames)
		if err != nil {
			return "", err
		}
		text += nextText
	}
	return text, nil
}

// GetBmarkTextDetailed returns detailed, multi-line bookmark text used for an ephemeral post
func (b *Bookmarks) GetBmarkTextDetailed(bmark *Bookmark, labelNames []string, args *model.CommandArgs) (string, error) {
	title, err := b.getTitleFromPost(bmark.PostID)
	if err != nil {
		return "", err
	}

	if bmark.HasUserTitle() {
		title = bmark.Title
	}

	codeBlockedNames := GetCodeBlockedLabels(labelNames)
	post, appErr := b.api.GetPost(bmark.PostID)
	if appErr != nil {
		return "", appErr
	}

	iconLink := getIconLink(b.api, bmark.PostID)

	text := fmt.Sprintf("%s\n#### Bookmark Title %s\n", codeBlockedNames, iconLink)
	text += fmt.Sprintf("**%s**\n", title)
	text += "##### Post Message \n"
	text += post.Message

	return text, nil
}
