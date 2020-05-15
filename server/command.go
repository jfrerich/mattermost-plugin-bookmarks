package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	commandTriggerBookmarks = "bookmarks"

	addCommandText = `
**/bookmarks add**
* |/bookmarks add <post_id> <bookmark_title> --labels <label1,label2>| - add a bookmark by specifying a post_id (with optional title)
* |/bookmarks add <permalink> <bookmark_title> --labels <label1,label2>| - add a bookmark by specifying the post permalink (with optional title)
`
	labelCommandText = `
**/bookmarks label**
* |/bookmarks label <post_id> --labels <labels>| - add labels (comma-separated) to a bookmark
* |/bookmarks label add <labels> | - create a new label
* |/bookmarks label remove <labels> | - remove a label
* |/bookmarks label remove <labels> --force | - forces removal of labels from
* bookmarks currently using the label as well as the label list
* |/bookmarks label view | - list all labels
`
	viewCommandText = `
**/bookmarks view**
* |/bookmarks view| - view all saved bookmarks
* |/bookmarks view <post_id> OR <permalink>| - view detailed bookmark view
`
	removeCommandText = `
**/bookmarks remove**
* |/bookmarks remove <post_id>| - remove bookmarks by post_id, or permalink
* |/bookmarks remove <post_id1> <post_id2>| - remove multiple bookmarks by post_id, or permalink
`
	// 	renameCommandText = `
	// **/bookmarks rename**
	// * |/bookmarks rename <label-old> <label-new>| - rename a label
	// `
	helpCommandText = `###### Bookmarks Slash Command Help` +
		addCommandText +
		labelCommandText +
		viewCommandText +
		removeCommandText
	// renameCommandText
)

func getHelp(text string) string {
	return strings.Replace(text, "|", "`", -1)
}

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          commandTriggerBookmarks,
		DisplayName:      commandTriggerBookmarks,
		Description:      "Manage Mattermost messages!",
		AutoComplete:     true,
		AutoCompleteHint: "[command]",
		AutoCompleteDesc: "Available commands: add, view, remove, label help",
	}
}

func (p *Plugin) postCommandResponse(args *model.CommandArgs, text string) {
	post := &model.Post{
		UserId:    p.getBotID(),
		ChannelId: args.ChannelId,
		Message:   text,
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)
}

func (p *Plugin) responsef(commandArgs *model.CommandArgs, format string, args ...interface{}) *model.CommandResponse {
	p.postCommandResponse(commandArgs, fmt.Sprintf(format, args...))
	return &model.CommandResponse{}
}

func (p *Plugin) executeCommandHelp(args *model.CommandArgs) *model.CommandResponse {
	return p.responsef(args, getHelp(helpCommandText))
}

// ExecuteCommand executes a command that has been previously registered via the RegisterCommand API.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	if len(split) < 2 {
		return p.executeCommandHelp(args), nil
	}

	action := split[1]

	//nolint:goconst
	switch action {
	case "add":
		return p.executeCommandAdd(args), nil
	case "label":
		return p.executeCommandLabel(args), nil
	case "remove":
		return p.executeCommandRemove(args), nil
	case "view":
		return p.executeCommandView(args), nil
	case "help":
		return p.executeCommandHelp(args), nil

	default:
		return p.responsef(args, fmt.Sprintf("Unknown command: "+args.Command)), nil
	}
}
