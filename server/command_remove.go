package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

// executeCommandRemove removes a given bookmark from the store
func (p *Plugin) executeCommandRemove(args *model.CommandArgs) *model.CommandResponse {
	subCommand := strings.Fields(args.Command)

	if len(subCommand) < 3 {
		return p.responsef(args, "Missing sub-command. You can try %v", getHelp(removeCommandText))
	}

	bookmarkIDs := subCommand[2:]

	text := "Removed bookmark: "
	if len(bookmarkIDs) > 1 {
		text = "Removed bookmarks: \n"
	}

	b := NewBookmarksWithUser(p.API, args.UserId)
	bmarks, err := b.getBookmarks()
	if err != nil {
		return p.responsef(args, err.Error())
	}
	if bmarks == nil {
		return p.responsef(args, "User doesn't have any bookmarks")
	}

	labels := NewLabelsWithUser(p.API, args.UserId)
	labels, err = labels.getLabelsForUser()
	if err != nil {
		return p.responsef(args, "Unable to get labels for user, %s", err)
	}

	for _, id := range bookmarkIDs {
		bookmarkID := p.getPostIDFromLink(id)
		bmark, err := bmarks.getBookmark(bookmarkID)
		if err != nil {
			return p.responsef(args, err.Error())
		}

		var labelNames []string
		for _, labelID := range bmark.LabelIDs {
			name, _ := labels.getNameFromID(labelID)
			labelNames = append(labelNames, name)
		}

		newText, err := p.getBmarkTextOneLine(bmark, labelNames, args)
		// FIXME newText, err := p.getBmarkTextOneLine(bmark, args)
		if err != nil {
			return p.responsef(args, err.Error())
		}

		_, err = bmarks.deleteBookmark(bookmarkID)
		if err != nil {
			return p.responsef(args, err.Error())
		}

		text += newText
	}

	return p.responsef(args, fmt.Sprint(text))
}
