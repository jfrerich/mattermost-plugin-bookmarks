package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

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
	// case "remove":
	// 	return p.executeCommandRemove(args), nil
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

	names := subCommand[3:]

	text := "Added Label: "
	if len(names) > 1 {
		text = "Added Labels: \n"
	}
	text = text + strings.Join(names, "\n")

	_, err := p.addLabels(args.UserId, names)
	if err != nil {
		return p.responsef(args, "Unable to add labels for user %s", args.UserId)
	}

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

	for _, label := range labels.ByName {
		v := fmt.Sprintf("`%s`\n", label.Name)
		text = text + v
	}

	return p.responsef(args, fmt.Sprintf(text))
}
