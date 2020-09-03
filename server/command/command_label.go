package command

import (
	"fmt"
	"strings"

	"github.com/jfrerich/mattermost-plugin-bookmarks/server/bookmarks"

	"github.com/spf13/pflag"
)

const (
	flagForce = "force"
)

type removeLabelOptions struct {
	force bool
}

func getLabelRemoveFlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("remove labels", pflag.ContinueOnError)
	flagSet.Bool(flagForce, false, "force removal of labels when they currently exist on a bookmark")

	return flagSet
}

func parseLabelRemoveArgs(args []string) (removeLabelOptions, error) {
	var options removeLabelOptions
	removeLabelFlagSet := getLabelRemoveFlagSet()
	err := removeLabelFlagSet.Parse(args)
	if err != nil {
		return options, err
	}

	options.force, err = removeLabelFlagSet.GetBool(flagForce)
	if err != nil {
		return options, err
	}

	return options, nil
}

// ExecuteCommandLabel executes a label sub-command
func (c *Command) executeCommandLabel() string {
	split := strings.Fields(c.Args.Command)
	if len(split) < 3 {
		return c.responsef(c.Args, "Missing label sub-command. You can try %v", getHelp(labelCommandText))
	}

	action := split[2]

	handler := c.responsef(c.Args, fmt.Sprintf("Unknown command: "+c.Args.Command))
	switch action {
	case "add":
		handler = c.executeCommandLabelAdd()
	case "remove":
		handler = c.executeCommandLabelRemove()
	case "rename":
		handler = c.executeCommandLabelRename()
	case "view":
		handler = c.executeCommandLabelView()
	case "help":
		handler = c.responsef(c.Args, "Please specify a label name %v", getHelp(labelCommandText))
	}
	return handler
}

func (c *Command) executeCommandLabelAdd() string {
	subCommand := strings.Fields(c.Args.Command)
	if len(subCommand) < 4 {
		return c.responsef(c.Args, "Please specify a label name %v", getHelp(labelCommandText))
	}

	labelName := subCommand[3]

	labels, err := bookmarks.NewLabelsWithUser(c.API, c.Args.UserId)
	if err != nil {
		return c.responsef(c.Args, err.Error())
	}

	_, err = labels.AddLabel(labelName)
	if err != nil {
		return c.responsef(c.Args, err.Error())
	}

	text := "Added Label: "
	text += fmt.Sprintf("%v", labelName)

	return c.responsef(c.Args, fmt.Sprint(text))
}

func (c *Command) executeCommandLabelRename() string {
	subCommand := strings.Fields(c.Args.Command)
	if len(subCommand) < 5 {
		return c.responsef(c.Args, "Please specify a `to` and `from` label name%v", getHelp(labelCommandText))
	}

	from := subCommand[3]
	to := subCommand[4]

	labels, err := bookmarks.NewLabelsWithUser(c.API, c.Args.UserId)
	if err != nil {
		return c.responsef(c.Args, err.Error())
	}

	lfrom := labels.GetLabelByName(from)
	if lfrom == nil {
		return c.responsef(c.Args, fmt.Sprintf("Label `%v` does not exist", from))
	}

	// if the "to" label already exists, alert the user with options
	lto := labels.GetLabelByName(to)
	if lto != nil {
		return c.responsef(c.Args, fmt.Sprintf("Cannot rename Label `%v` to `%v`. Label already exists. Please choose a different label name", from, to))
	}

	lfrom.Name = to
	labels.ByID[lfrom.ID] = lfrom
	if err := labels.StoreLabels(); err != nil {
		return c.responsef(c.Args, "failed to add label")
	}

	text := fmt.Sprintf("Renamed label from `%v` to `%v`", from, to)

	return c.responsef(c.Args, fmt.Sprint(text))
}

// executeCommandLabelRemove removes a given bookmark from the store
func (c *Command) executeCommandLabelRemove() string {
	subCommand := strings.Fields(c.Args.Command)
	if len(subCommand) < 4 {
		return c.responsef(c.Args, "Please specify a label name %v", getHelp(labelCommandText))
	}

	labels, err := bookmarks.NewLabelsWithUser(c.API, c.Args.UserId)
	if err != nil {
		return c.responsef(c.Args, err.Error())
	}
	if labels == nil || len(labels.ByID) == 0 {
		return c.responsef(c.Args, "You do not have any saved labels")
	}

	labelName := subCommand[3]

	labelID, err := labels.GetIDFromName(labelName)
	if err != nil {
		return c.responsef(c.Args, err.Error())
	}
	bmarks, err := bookmarks.NewBookmarksWithUser(c.API, c.Args.UserId)
	if err != nil {
		return c.responsef(c.Args, err.Error())
	}

	options, err := parseLabelRemoveArgs(subCommand)
	if err != nil {
		return c.responsef(c.Args, "Unable to parse options, %s", err)
	}

	if bmarks != nil {
		// check to see if any bookmarks currently have the label
		bmarks, err = bmarks.GetBookmarksWithLabelID(labelID)
		if err != nil {
			return c.responsef(c.Args, err.Error())
		}

		numBmarksWithLabel := len(bmarks.ByID)
		if numBmarksWithLabel != 0 && !options.force {
			return c.responsef(
				c.Args,
				fmt.Sprintf("There are %v bookmarks with the label:%s. Use the `--force` flag remove the label from the bookmarks.",
					numBmarksWithLabel, bookmarks.GetCodeBlockedLabels([]string{labelName})),
			)
		}

		// delete label from bookmarks
		for _, bmark := range bmarks.ByID {
			err = bmarks.DeleteLabel(bmark.PostID, labelID)
			if err != nil {
				return c.responsef(c.Args, err.Error())
			}
		}
	}

	// delete from store after delete from bookmarks
	err = labels.DeleteByID(labelID)
	if err != nil {
		return c.responsef(c.Args, err.Error())
	}

	text := "Removed label: "
	text += fmt.Sprintf("`%v`", labelName)
	return c.responsef(c.Args, fmt.Sprint(text))
}

func (c *Command) executeCommandLabelView() string {
	subCommand := strings.Fields(c.Args.Command)

	if len(subCommand) != 3 {
		return c.responsef(c.Args, "view subcommand takes no arguments%v", getHelp(labelCommandText))
	}

	labels, err := bookmarks.NewLabelsWithUser(c.API, c.Args.UserId)
	if err != nil {
		return c.responsef(c.Args, err.Error())
	}
	if labels == nil || len(labels.ByID) == 0 {
		return c.responsef(c.Args, "You do not have any saved labels")
	}

	text := "#### Labels List\n"
	for _, label := range labels.ByID {
		v := fmt.Sprintf("`%s`\n", label.Name)
		text += v
	}

	return c.responsef(c.Args, fmt.Sprint(text))
}
