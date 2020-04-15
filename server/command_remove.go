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

	bookmarkIDs := strings.Split(subCommand[2], ",")

	text := "Removed bookmark: "
	if len(bookmarkIDs) > 1 {
		text = "Removed bookmarks: \n"
	}

	for _, id := range bookmarkIDs {
		bookmarkID := p.getPostIDFromLink(id)

		bmark, err := p.deleteBookmark(args.UserId, bookmarkID)
		if err != nil {
			return p.responsef(args, err.Error())
		}

		newText, appErr := p.getBmarkTextOneLine(bmark, args)
		if appErr != nil {
			return p.responsef(args, "Unable to get bookmarks list bookmark")
		}
		text = text + newText
	}

	return p.responsef(args, fmt.Sprintf(text))
}
