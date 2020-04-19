package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
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
func (p *Plugin) executeCommandLabel(args *model.CommandArgs) *model.CommandResponse {

	split := strings.Fields(args.Command)
	if len(split) < 3 {
		return p.responsef(args, "Missing label sub-command. You can try %v", getHelp(labelCommandText))
	}

	action := split[2]

	switch action {
	case "add":
		return p.executeCommandLabelAdd(args)
	case "remove":
		return p.executeCommandLabelRemove(args)
	// case "rename":
	// 	return p.executeCommandLabelRename(args)
	case "view":
		return p.executeCommandLabelView(args)
	case "help":
		return p.responsef(args, "Please specify a label name %v", getHelp(labelCommandText))

	default:
		return p.responsef(args, fmt.Sprintf("Unknown command: "+args.Command))
	}
}

func (p *Plugin) executeCommandLabelAdd(args *model.CommandArgs) *model.CommandResponse {
	subCommand := strings.Fields(args.Command)
	if len(subCommand) < 4 {
		return p.responsef(args, "Please specify a label name %v", getHelp(labelCommandText))
	}

	labelName := subCommand[3]

	l := NewLabels(p.API)
	labels, err := l.getLabels(args.UserId)

	label, err := labels.addLabel(args.UserId, labelName)
	if err != nil {
		return p.responsef(args, err.Error())
	}

	text := "Added Label: "
	text = text + fmt.Sprintf("%v", label.Name)

	return p.responsef(args, fmt.Sprintf(text))
}

// executeCommandLabelRemove removes a given bookmark from the store
func (p *Plugin) executeCommandLabelRemove(args *model.CommandArgs) *model.CommandResponse {
	subCommand := strings.Fields(args.Command)

	if len(subCommand) < 4 {
		return p.responsef(args, "Please specify a label name %v", getHelp(labelCommandText))
	}

	l := NewLabels(p.API)
	labels, err := l.getLabels(args.UserId)
	if err != nil {
		return p.responsef(args, err.Error())
	}
	if labels == nil || len(labels.ByID) == 0 {
		return p.responsef(args, "You do not have any saved labels")
	}

	labelName := subCommand[3]

	labelID, err := labels.getIDFromName(args.UserId, labelName)
	if err != nil {
		return p.responsef(args, err.Error())
	}
	b := NewBookmarks(p.API)
	bmarks, err := b.getBookmarks(args.UserId)
	if err != nil {
		return p.responsef(args, "Unable to retrieve bookmarks for user %s", args.UserId)
	}

	options, err := parseLabelRemoveArgs(subCommand)
	if err != nil {
		return p.responsef(args, "Unable to parse options, %s", err)
	}

	if bmarks != nil {
		// check to see if any bookmarks currently have the label
		bmarks, err = bmarks.getBookmarksWithLabelID(args.UserId, labelID)
		if err != nil {
			return p.responsef(args, err.Error())
		}
		numBmarksWithLabel := len(bmarks.ByID)
		if numBmarksWithLabel != 0 && !options.force {
			return p.responsef(
				args,
				fmt.Sprintf("There are %v bookmarks with the label:%s. Use the --force flag remove the label from the bookmarks.",
					numBmarksWithLabel, p.getCodeBlockedLabels([]string{labelName})),
			)
		}

		// delete label from bookmarks
		for _, bmark := range bmarks.ByID {
			err = bmarks.deleteLabel(args.UserId, bmark.PostID, labelID)
			if err != nil {
				return p.responsef(args, err.Error())
			}
		}
	}

	// delete from store after delete from bookmarks
	err = labels.deleteByID(args.UserId, labelID)
	if err != nil {
		return p.responsef(args, err.Error())
	}

	text := "Removed label: "
	text = text + fmt.Sprintf("`%v`", labelName)
	return p.responsef(args, fmt.Sprintf(text))
}

func (p *Plugin) executeCommandLabelView(args *model.CommandArgs) *model.CommandResponse {
	subCommand := strings.Fields(args.Command)

	if len(subCommand) != 3 {
		return p.responsef(args, "view subcommand takes no arguments%v", getHelp(labelCommandText))
	}

	l := NewLabels(p.API)
	labels, err := l.getLabels(args.UserId)
	if err != nil {
		return p.responsef(args, "Unable to retrieve bookmark for user %s", args.UserId)
	}
	if labels == nil || len(labels.ByID) == 0 {
		return p.responsef(args, "You do not have any saved labels")
	}

	text := "#### Labels List\n"
	for _, label := range labels.ByID {
		v := fmt.Sprintf("`%s`\n", label.Name)
		text = text + v
	}

	return p.responsef(args, fmt.Sprintf(text))
}
