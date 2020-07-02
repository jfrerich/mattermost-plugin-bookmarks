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
* |/bookmarks label rename <old> <new>| - rename a label
* |/bookmarks label remove <labels> | - remove a label
* |/bookmarks label remove <labels> --force | - forces removal of labels from bookmarks currently using the label as well as the label list
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
	helpCommandText = `###### Bookmarks Slash Command Help` +
		addCommandText +
		labelCommandText +
		viewCommandText +
		removeCommandText
)

func getHelp(text string) string {
	return strings.Replace(text, "|", "`", -1)
}

func createBookmarksCommand() *model.Command {
	bookmarks := model.NewAutocompleteData(
		commandTriggerBookmarks, "[command]", "Available commands: add, label, remove, view, help")

	// top-level commands
	bookmarks.AddCommand(createAddCommand())
	bookmarks.AddCommand(createLabelCommand())
	bookmarks.AddCommand(createRemoveCommand())
	bookmarks.AddCommand(createViewCommand())
	bookmarks.AddCommand(createHelpCommand())

	return &model.Command{
		Trigger:          commandTriggerBookmarks,
		DisplayName:      commandTriggerBookmarks,
		Description:      "Manage Mattermost messages!",
		AutoComplete:     true,
		AutocompleteData: bookmarks,
		AutoCompleteHint: "[command]",
		AutoCompleteDesc: "Available commands: add, label, remove, view, help",
	}
}

func prefixWithAPI(route string) string {
	return routeAPIPrefix + route
}

// createHelpCommand adds the help autocomplete option
func createHelpCommand() *model.AutocompleteData {
	add := model.NewAutocompleteData(
		"help", "", "show help")
	return add
}

// createAddCommand adds the add autocomplete option
func createAddCommand() *model.AutocompleteData {
	add := model.NewAutocompleteData(
		"add", "[post-id OR permalink] --labels", "Add a bookmark")
	return add
}

// createLabelCommand adds the label autocomplete with suboptions
func createLabelCommand() *model.AutocompleteData {
	label := model.NewAutocompleteData(
		"label", "[add|remove|rename|view]", "Create, remove, modify, or view labels")
	label.AddCommand(createLabelAddCommand())
	label.AddCommand(createLabelRemoveCommand())
	label.AddCommand(createLabelRenameCommand())
	label.AddCommand(createLabelViewCommand())
	return label
}

func createLabelRemoveCommand() *model.AutocompleteData {
	remove := model.NewAutocompleteData(
		"remove", "[label-name] --force", "Remove a label")
	remove.AddDynamicListArgument("Label Name", prefixWithAPI(routeAutocompleteLabels), false)
	return remove
}

func createLabelViewCommand() *model.AutocompleteData {
	remove := model.NewAutocompleteData(
		"view", "", "View all labels")
	return remove
}

func createLabelAddCommand() *model.AutocompleteData {
	remove := model.NewAutocompleteData(
		"add", "[label-name]", "add a label")
	remove.AddDynamicListArgument("Label Name", "", false)
	return remove
}

func createLabelRenameCommand() *model.AutocompleteData {
	remove := model.NewAutocompleteData(
		"rename", "[label-name] [new-label-name]", "Rename a label")
	remove.AddDynamicListArgument("Label Name", prefixWithAPI(routeAutocompleteLabels), false)
	return remove
}

// createRemoveCommand adds the remove autocomplete option
func createRemoveCommand() *model.AutocompleteData {
	remove := model.NewAutocompleteData(
		"remove", "[post-id]", "Remove a bookmark")
	remove.AddDynamicListArgument("[post_id] OR [permalink]", prefixWithAPI(routeAutocompleteBookmarks), false)
	return remove
}

// createViewCommand adds the View autocomplete option with suboptions
func createViewCommand() *model.AutocompleteData {
	view := model.NewAutocompleteData(
		"view", "[post_id] OR [permalink]", "View a bookmark or all bookmarks")
	view.AddDynamicListArgument("[post_id] OR [permalink]", prefixWithAPI(routeAutocompleteBookmarks), false)
	return view
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
