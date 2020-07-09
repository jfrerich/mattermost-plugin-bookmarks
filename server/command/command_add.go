package command

import (
	"strings"

	"github.com/jfrerich/mattermost-plugin-bookmarks/server/bookmarks"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/utils"

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
func (c *Command) executeCommandAdd() string {
	subCommand := strings.Fields(c.Args.Command)
	subCommand = subCommand[2:]

	if len(subCommand) < 1 {
		return c.responsef(c.Args, "Missing sub-command. You can try %v", getHelp(addCommandText))
	}
	postID := utils.GetPostIDFromLink(subCommand[0])

	_, appErr := c.API.GetPost(postID)
	if appErr != nil {
		return c.responsef(c.Args, "PostID `%s` is not a valid postID", postID)
	}

	var bookmark bookmarks.Bookmark
	bookmark.PostID = postID

	// user provides a title
	if len(subCommand) >= 2 {
		title := c.getTitleFromArguments(subCommand[1:])
		bookmark.SetTitle(title)
	}

	options, err := parseAddBookmarkArgs(subCommand)
	if err != nil {
		return c.responsef(c.Args, "Unable to parse options, %s", err)
	}

	var labelIDsForBookmark []string

	// user going to add labels names
	if len(options.labels) != 0 {
		var labels *bookmarks.Labels
		labels, err = bookmarks.NewLabelsWithUser(c.API, c.Args.UserId)
		if err != nil {
			return c.responsef(c.Args, "Unable to get labels for user, %s", err)
		}

		for _, name := range options.labels {
			label := labels.GetLabelByName(name)
			// create new label in labels store and add ID to bookmark
			if label == nil {
				_, err = labels.AddLabel(name)
				if err != nil {
					return c.responsef(c.Args, "Unable to add new label for: %s, err=%s", name, err.Error())
				}
			}
			var labelID string
			labelID, err = labels.GetIDFromName(name)
			if err != nil {
				return c.responsef(c.Args, err.Error())
			}
			labelIDsForBookmark = append(labelIDsForBookmark, labelID)
		}
		bookmark.AddLabelIDs(labelIDsForBookmark)
	}

	// get all bookmarks for user
	bmarks, err := bookmarks.NewBookmarksWithUser(c.API, c.Args.UserId)
	if err != nil {
		return c.responsef(c.Args, "Unable to get bookmarks")
	}

	err = bmarks.AddBookmark(&bookmark)
	if err != nil {
		return c.responsef(c.Args, "Unable to add bookmark")
	}

	text, err := bmarks.GetBmarkTextOneLine(&bookmark, options.labels)
	if err != nil {
		return c.responsef(c.Args, "Unable to get bookmarks list bookmark")
	}

	return c.responsef(c.Args, "Added bookmark: %s", text)
}

func (c *Command) getTitleFromArguments(args []string) string {
	for i, arg := range args {
		// user also provided a --flag
		if strings.HasPrefix(arg, "--") {
			title := strings.Join(args[:i], " ")
			return title
		}
	}

	// user provided no flags after the ID, rejoin with spaces and return
	return strings.Join(args, " ")
}
