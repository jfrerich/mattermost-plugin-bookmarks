package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

// executeCommandView shows all bookmarks in an ephemeral post
func (p *Plugin) executeCommandView(args *model.CommandArgs) *model.CommandResponse {
	subCommand := strings.Fields(args.Command)

	b := NewBookmarksWithUser(p.API, args.UserId)
	bmarks, err := b.getBookmarks()
	if err != nil {
		return p.responsef(args, "Unable to retrieve bookmarks for user %s", args.UserId)
	}

	// bookmarks is nil if user has never added a bookmark.
	// bookmarks.ByID will be empty if user created a bookmark and then deleted
	// it and now has 0 bookmarks
	if bmarks == nil || len(bmarks.ByID) == 0 {
		return p.responsef(args, "You do not have any saved bookmarks")
	}

	labels := NewLabelsWithUser(p.API, args.UserId)
	labels, err = labels.getLabelsForUser()
	if err != nil {
		return p.responsef(args, "Unable to get labels for user, %s", err)
	}

	// user requests to view an individual bookmark
	if len(subCommand) == 3 {
		postID := subCommand[2]
		postID = p.getPostIDFromLink(postID)

		var bmark *Bookmark
		bmark, err = bmarks.getBookmark(postID)
		if err != nil {
			return p.responsef(args, err.Error())
		}

		var labelNames []string
		labelNames, err = labels.getLabelNames(args.UserId, bmark)
		if err != nil {
			return p.responsef(args, err.Error())
		}

		var text string
		text, err = p.getBmarkTextDetailed(bmark, labelNames, args)
		if err != nil {
			return p.responsef(args, "Unable to get bookmark text %s", err)
		}
		return p.responsef(args, text)
	}

	text := "#### Bookmarks List\n"
	bmarksSorted, err := b.ByPostCreateAt(bmarks)
	if err != nil {
		return p.responsef(args, err.Error())
	}

	for _, bmark := range bmarksSorted {
		labelNames, err := labels.getLabelNames(args.UserId, bmark)
		if err != nil {
			return p.responsef(args, err.Error())
		}
		nextText, err := p.getBmarkTextOneLine(bmark, labelNames, args)
		if err != nil {
			return p.responsef(args, "Unable to get bookmarks text %s", err)
		}
		text += nextText
	}

	return p.responsef(args, text)
}
