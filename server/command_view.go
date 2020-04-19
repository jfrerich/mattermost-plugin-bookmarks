package main

import (
	"fmt"
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
	if bmarks == nil {
		return p.responsef(args, fmt.Sprintf("You do not have any saved bookmarks"))
	}

	// user requests to view an indiviual bookmark
	if len(subCommand) == 3 {
		postID := subCommand[2]
		postID = p.getPostIDFromLink(postID)
		bmark, err := bmarks.getBookmark(args.UserId, postID)
		if err != nil {
			return p.responsef(args, "Unable to retrieve bookmark for user %s", args.UserId)
		}

		text, appErr := p.getBmarkTextDetailed(bmark, args)
		if appErr != nil {
			return p.responsef(args, "Unable to retrieve bookmark for user %s", args.UserId)
		}
		return p.responsef(args, text)
	}

	// bookmarks is nil if user has never added a bookmark.
	// bookmarks.ByID will be empty if user created a bookmark and then deleted
	// it and now has 0 bookmarks
	if bmarks == nil || len(bmarks.ByID) == 0 {
		return p.responsef(args, "You do not have any saved bookmarks")
	}

	text := "#### Bookmarks List\n"
	bmarksSorted, appErr := b.ByPostCreateAt(bmarks)
	if appErr != nil {
		return p.responsef(args, "Unable to retrieve bookmarks for user %s", args.UserId)
	}

	for _, bmark := range bmarksSorted {
		nextText, appErr := p.getBmarkTextOneLine(bmark, args)
		if appErr != nil {
			return p.responsef(args, "Unable to get bookmarks list bookmark")
		}
		text = text + nextText
	}

	return p.responsef(args, text)
}
