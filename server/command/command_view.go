package command

import (
	"strings"

	"github.com/jfrerich/mattermost-plugin-bookmarks/server/bookmarks"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/utils"
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
func (c *Command) executeCommandView() string {
	subCommand := strings.Fields(c.Args.Command)

	bmarks, err := bookmarks.NewBookmarksWithUser(c.API, c.Args.UserId)
	if err != nil {
		return c.responsef(c.Args, "Unable to retrieve bookmarks for user %s", c.Args.UserId)
	}

	// bookmarks is nil if user has never added a bookmark.
	// bookmarks.ByID will be empty if user created a bookmark and then deleted
	// it and now has 0 bookmarks
	if bmarks == nil || len(bmarks.ByID) == 0 {
		return c.responsef(c.Args, "You do not have any saved bookmarks")
	}

	// user requests to view an individual bookmark
	if len(subCommand) == 3 && !strings.HasPrefix(subCommand[2], "--") {
		postID := subCommand[2]
		postID = utils.GetPostIDFromLink(postID)
		text, _ := c.commandViewPostID(postID, bmarks)
		return c.responsef(c.Args, text)
	}

	options, err := parseViewBookmarkArgs(subCommand)
	if err != nil {
		return c.responsef(c.Args, "Unable to parse options, %s", err)
	}

	var bmarkFilters bookmarks.Filters
	bmarkFilters.LabelNames = options.labels

	text, err := bmarks.GetBmarksEphemeralText(c.Args.UserId, &bmarkFilters)
	if err != nil {
		return c.responsef(c.Args, text)
	}

	return c.responsef(c.Args, text)
}

// executeCommandView shows all bookmarks in an ephemeral post
func (c *Command) commandViewPostID(postID string, bmarks *bookmarks.Bookmarks) (string, error) {
	postID = utils.GetPostIDFromLink(postID)

	var bmark *bookmarks.Bookmark
	bmark, err := bmarks.GetBookmark(postID)
	if err != nil {
		return "", err
	}

	var labelNames []string
	labelNames, err = bmarks.GetBmarkLabelNames(bmark)
	if err != nil {
		return "", err
	}

	var text string
	text, err = bmarks.GetBmarkTextDetailed(bmark, labelNames, c.Args)
	if err != nil {
		return "", errors.Wrap(err, "Unable to get bookmark text")
	}
	return text, nil
}
