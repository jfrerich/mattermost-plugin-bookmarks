package command

import (
	"fmt"
	"strings"

	"github.com/jfrerich/mattermost-plugin-bookmarks/server/bookmarks"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/utils"
)

// executeCommandRemove removes a given bookmark from the store
func (c *Command) executeCommandRemove() string {
	subCommand := strings.Fields(c.Args.Command)

	if len(subCommand) < 3 {
		return c.responsef(c.Args, "Missing sub-command. You can try %v", getHelp(removeCommandText))
	}

	bookmarkIDs := subCommand[2:]

	text := "Removed bookmark: "
	if len(bookmarkIDs) > 1 {
		text = "Removed bookmarks: \n"
	}

	bmarks, err := bookmarks.NewBookmarksWithUser(c.API, c.Args.UserId)
	if err != nil {
		return c.responsef(c.Args, err.Error())
	}
	if bmarks == nil {
		return c.responsef(c.Args, "User doesn't have any bookmarks")
	}

	labels, err := bookmarks.NewLabelsWithUser(c.API, c.Args.UserId)
	if err != nil {
		return c.responsef(c.Args, "Unable to get labels for user, %s", err)
	}

	for _, id := range bookmarkIDs {
		bookmarkID := utils.GetPostIDFromLink(id)
		bmark, err := bmarks.GetBookmark(bookmarkID)
		if err != nil {
			return c.responsef(c.Args, err.Error())
		}

		var labelNames []string
		for _, labelID := range bmark.LabelIDs {
			name, _ := labels.GetNameFromID(labelID)
			labelNames = append(labelNames, name)
		}

		newText, err := bmarks.GetBmarkTextOneLine(bmark, labelNames)
		if err != nil {
			return c.responsef(c.Args, err.Error())
		}

		err = bmarks.DeleteBookmark(bookmarkID)
		if err != nil {
			return c.responsef(c.Args, err.Error())
		}

		text += newText
	}

	return c.responsef(c.Args, fmt.Sprint(text))
}
