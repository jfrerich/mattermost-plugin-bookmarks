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
		title := p.getTitleFromArguments(subCommand[1:])
		bookmark.setTitle(title)
	}

	options, err := parseAddBookmarkArgs(subCommand)
	if err != nil {
		return p.responsef(args, "Unable to parse options, %s", err)
	}

	// var labels *Labels
	var labelIDsForBookmark []string

	// user going to add labels names
	if len(options.labels) != 0 {

		labels := NewLabelsWithUser(p.API, args.UserId)
		labels, err = labels.getLabels()
		if err != nil {
			return p.responsef(args, "Unable to get labels for user, %s", err)
		}

		for _, name := range options.labels {
			label := labels.getLabelByName(name)
			// create new label in labels store and add ID to bookmark
			if label == nil {
				_, err = labels.addLabel(name)
				if err != nil {
					return p.responsef(args, "Unable to add new label for: %s, err=%s", name, err.Error())
				}
			}
			var labelID string
			labelID, err = labels.getIDFromName(name)
			if err != nil {
				return p.responsef(args, err.Error())
			}
			labelIDsForBookmark = append(labelIDsForBookmark, labelID)
		}
		bookmark.addLabelIDs(labelIDsForBookmark)
	}

	// get all bookmarks for user
	b := NewBookmarksWithUser(p.API, args.UserId)
	bmarks, err := b.getBookmarks()
	if err != nil {
		return p.responsef(args, "Unable to get bookmarks")
	}

	// no marks, initialize the store first
	if bmarks == nil {
		bmarks = NewBookmarksWithUser(p.API, args.UserId)
	}
	bmarks.addBookmark(&bookmark)

	text, err := p.getBmarkTextOneLine(&bookmark, options.labels, args)
	if err != nil {
		return p.responsef(args, "Unable to get bookmarks list bookmark")
	}

	return p.responsef(args, "Added bookmark: %s", text)
}

func (p *Plugin) getTitleFromArguments(args []string) string {
	for i, arg := range args {
		// user also provided a --flag
		if arg == "--" {
			return strings.Join(args[:i-1], " ")
		}
	}

	// user provided no flags after the ID, rejoin with spaces and return
	return strings.Join(args, " ")
}
