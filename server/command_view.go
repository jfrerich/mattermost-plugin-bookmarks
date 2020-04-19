package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/spf13/pflag"
)

func getViewBookmarksFlagSet() *pflag.FlagSet {
	getViewBookmarksFlagSet := pflag.NewFlagSet("view bookmarks", pflag.ContinueOnError)
	getViewBookmarksFlagSet.String("filter-labels", "", "Filter bookmarks with these specified labels")

	return getViewBookmarksFlagSet
}

// executeCommandView shows all bookmarks in an ephemeral post
func (p *Plugin) executeCommandView(args *model.CommandArgs) *model.CommandResponse {

	subCommand := strings.Fields(args.Command)

	b := NewBookmarks(p.API)
	bmarks, err := b.getBookmarks(args.UserId)
	if err != nil {
		return p.responsef(args, "Unable to retrieve bookmarks for user %s", args.UserId)
	}

	// bookmarks is nil if user has never added a bookmark.
	// bookmarks.ByID will be empty if user created a bookmark and then deleted
	// it and now has 0 bookmarks
	if bmarks == nil || len(bmarks.ByID) == 0 {
		return p.responsef(args, "You do not have any saved bookmarks")
	}

	// user requests to view an individual bookmark
	if len(subCommand) == 3 {
		postID := subCommand[2]
		postID = p.getPostIDFromLink(postID)
		bmark, err := bmarks.getBookmark(args.UserId, postID)
		if err != nil {
			return p.responsef(args, err.Error())
		}

		text, err := p.getBmarkTextDetailed(bmark, args)
		if err != nil {
			return p.responsef(args, err.Error())
		}
		return p.responsef(args, text)
	}

	text := "#### Bookmarks List\n"
	bmarksSorted, err := b.ByPostCreateAt(bmarks)
	if err != nil {
		return p.responsef(args, "Unable to retrieve bookmarks for user %s", args.UserId)
	}

	for _, bmark := range bmarksSorted {
		nextText, err := p.getBmarkTextOneLine(bmark, args)
		if err != nil {
			return p.responsef(args, "Unable to get bookmarks list bookmark")
		}
		text = text + nextText
	}

	return p.responsef(args, text)
}
