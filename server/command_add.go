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
		bookmark.Title = title
	}

	options, err := parseAddBookmarkArgs(subCommand)
	if err != nil {
		return p.responsef(args, "Unable to parse options, %s", err)
	}

	if len(options.labels) != 0 {
		bookmark.LabelNames = options.labels
	}

	p.addBookmark(args.UserId, &bookmark)

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
