package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/spf13/pflag"
)

const (
	flagLabel = "labels"
)

type addBookmarkOptions struct {
	labels []string
}

func getAddBookmarkFlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("add labels to bookmarks", pflag.ContinueOnError)
	flagSet.StringSlice(flagLabel, nil, "Add a label to a bookmark")

	return flagSet
}

func parseAddBookmarkArgs(args []string) (addBookmarkOptions, error) {
	var options addBookmarkOptions

	addBookmarkFlagSet := getAddBookmarkFlagSet()
	err := addBookmarkFlagSet.Parse(args)
	if err != nil {
		return options, err
	}

	options.labels, err = addBookmarkFlagSet.GetStringSlice(flagLabel)
	if err != nil {
		return options, err
	}

	return options, nil
}

// executeCommandAdd adds a bookmark to the store
func (p *Plugin) executeCommandAdd(args *model.CommandArgs) *model.CommandResponse {

	subCommand := strings.Fields(args.Command)
	subCommand = subCommand[2:]

	if len(subCommand) < 1 {
		return p.responsef(args, "Missing sub-command. You can try %v", getHelp(addCommandText))
	}
	postID := p.getPostIDFromLink(subCommand[0])

	_, appErr := p.API.GetPost(postID)
	if appErr != nil {
		return p.responsef(args, "PostID `%s` is not a valid postID", postID)
	}

	var bookmark Bookmark
	bookmark.PostID = postID

	// user provides a title
	if len(subCommand) >= 2 {
		title := p.constructValueFromArguments(subCommand[1:])
		bookmark.setTitle(title)
	}

	options, err := parseAddBookmarkArgs(subCommand)
	if err != nil {
		return p.responsef(args, "Unable to parse options, %s", err)
	}

	// user going to add labels names
	if len(options.labels) != 0 {

		l := NewLabels(p.API)
		labels, _ := l.getLabels(args.UserId)

		// get labelIDs from provided names
		// TODO: creates new label in labels table if name not found
		var labelUUIDs []string
		labelUUIDs, err = labels.getIDsFromNames(args.UserId, options.labels)
		if err != nil {
			return p.responsef(args, "Unable to get UUIDs for labels: %s", options.labels)
		}

		// add labelIDs to bmark
		bookmark.setLabelIDs(labelUUIDs)
	}

	// get all bookmarks for user
	b := NewBookmarks(p.API)
	bmarks, err := b.getBookmarks(args.UserId)
	if err != nil {
		return p.responsef(args, "Unable to get bookmarks")
	}

	// no marks, initialize the store first
	if bmarks == nil {
		bmarks = NewBookmarks(p.API)
	}

	bmarks.addBookmark(args.UserId, &bookmark)

	text, appErr := p.getBmarkTextOneLine(&bookmark, args)
	if appErr != nil {
		return p.responsef(args, "Unable to get bookmarks list bookmark")
	}

	return p.responsef(args, "Added bookmark: %s", text)
}

func (p *Plugin) constructValueFromArguments(args []string) string {
	index := 0
	for i, e := range args {
		// user also provided a --flag
		if e == "--" {
			return strings.Join(args[:index-1], " ")
		}
		index = i
	}

	// user provided no flags after the ID, rejoin with spaces and return
	return strings.Join(args, " ")
}
