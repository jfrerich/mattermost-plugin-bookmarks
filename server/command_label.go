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
	label, err := p.addLabel(args.UserId, labelName)
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

	labelName := subCommand[3]

	// check to see if any bookmarks currently have the label
	bmarks, err := p.getBookmarksWithLabel(args.UserId, labelName)
	if err != nil {
		return p.responsef(args, err.Error())
	}

	options, err := parseLabelRemoveArgs(subCommand)
	if err != nil {
		return p.responsef(args, "Unable to parse options, %s", err)
	}

	if bmarks != nil {
		numBmarksWithLabel := len(bmarks.ByID)
		if numBmarksWithLabel != 0 && !options.force {
			return p.responsef(
				args,
				fmt.Sprintf("There are %v bookmarks with the label:%s. Use the --force flag remove the label from the bookmarks.",
					numBmarksWithLabel, p.getCodeBlockedLabels([]string{labelName})),
			)
		}
	}

	label, err := p.deleteLabelByName(args.UserId, labelName)
	if err != nil {
		return p.responsef(args, err.Error())
	}
	if label == nil {
		return p.responsef(args, fmt.Sprintf("User doesn't have any labels"))
	}

	text := "Removed label: "
	text = text + fmt.Sprintf("`%v`", label.Name)
	return p.responsef(args, fmt.Sprintf(text))
}

func (p *Plugin) executeCommandLabelView(args *model.CommandArgs) *model.CommandResponse {
	subCommand := strings.Fields(args.Command)

	if len(subCommand) != 3 {
		return p.responsef(args, "view subcommand takes no arguments%v", getHelp(labelCommandText))
	}

	text := "#### Labels List\n"
	labels, err := p.getLabels(args.UserId)
	if err != nil {
		return p.responsef(args, "Unable to retrieve bookmark for user %s", args.UserId)
	}
	if labels == nil {
		return p.responsef(args, "You do not have any saved labels")
	}

	for _, label := range labels.ByID {
		v := fmt.Sprintf("`%s`\n", label.Name)
		text = text + v
	}

	return p.responsef(args, fmt.Sprintf(text))
}
