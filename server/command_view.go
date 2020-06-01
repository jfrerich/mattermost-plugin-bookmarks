package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const (
	flagFilterLabels = "filter-labels"
)

func getViewBookmarkFlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("filter bookmarks by label", pflag.ContinueOnError)
	flagSet.StringSlice(flagFilterLabels, nil, "filter by label")

	return flagSet
}

type viewBookmarkOptions struct {
	labels []string
}

func parseViewBookmarkArgs(args []string) (viewBookmarkOptions, error) {
	var options viewBookmarkOptions

	viewBookmarkFlagSet := getViewBookmarkFlagSet()
	err := viewBookmarkFlagSet.Parse(args)
	if err != nil {
		return options, err
	}

	options.labels, err = viewBookmarkFlagSet.GetStringSlice(flagFilterLabels)
	if err != nil {
		return options, err
	}

	return options, nil
}

// executeCommandView shows all bookmarks in an ephemeral post
func (p *Plugin) executeCommandView(args *model.CommandArgs) *model.CommandResponse {
	subCommand := strings.Fields(args.Command)

	bmarks, err := NewBookmarksWithUser(p.API, args.UserId).getBookmarks()
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
	if len(subCommand) == 3 && !strings.HasPrefix(subCommand[2], "--") {
		postID := subCommand[2]
		postID = p.getPostIDFromLink(postID)
		text, _ := p.commandViewPostID(postID, bmarks, args)
		return p.responsef(args, text)
	}

	options, err := parseViewBookmarkArgs(subCommand)
	if err != nil {
		return p.responsef(args, "Unable to parse options, %s", err)
	}

	var bmarkFilters BookmarksFilters
	bmarkFilters.LabelNames = options.labels

	text, err := p.getBmarksEphemeralText(args.UserId, &bmarkFilters)
	if err != nil {
		return p.responsef(args, text)
	}

	return p.responsef(args, text)
}

// executeCommandView shows all bookmarks in an ephemeral post
func (p *Plugin) commandViewPostID(postID string, bmarks *Bookmarks, args *model.CommandArgs) (string, error) {
	postID = p.getPostIDFromLink(postID)

	var bmark *Bookmark
	bmark, err := bmarks.getBookmark(postID)
	if err != nil {
		return "", err
	}

	var labelNames []string
	labelNames, err = bmarks.getBmarkLabelNames(bmark)
	if err != nil {
		return "", err
	}

	var text string
	text, err = p.getBmarkTextDetailed(bmark, labelNames, args)
	if err != nil {
		return "", errors.Wrap(err, "Unable to get bookmark text")
	}
	return text, nil
}
