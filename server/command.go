package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	// MaxTitleCharacters is the maximum length of characters displayed in a
	// bookmark title
	MaxTitleCharacters = 30

	commandTriggerBookmarks = "bookmarks"

	addCommandText = `
**/bookmarks add**
* |/bookmarks add <post_id OR post_permalink>| - bookmark a post_id with optional labels. if labels omitted, |unlabeled| autoadded
`
	labelCommandText = `
**/bookmarks label**
* |/bookmarks label <post_id> <labels>| - add labels to a bookmark; if labels omitted, |unlabeled| autoadded
* |/bookmarks label add <labels> | - create a new label
* |/bookmarks label list | - list all labels (include number of bookmarks per label)
`
	viewCommandText = `
**/bookmarks view**
* |/bookmarks view| - view bookmarks
`
	removeCommandText = `
**/bookmarks remove**
* |/bookmarks remove <post_id>| - remove labels from bookmarked post_id. if labels omitted remove post_id from bookmarks
`
	renameCommandText = `
**/bookmarks rename**
* |/bookmarks rename <label-old> <label-new>| - rename a label
`
	helpCommandText = `###### Bookmarks Slash Command Help` +
		addCommandText +
		// labelCommandText +
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
		AutoCompleteDesc: "Available commands: add, view, remove, help",
	}
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Username:     "bookmarks",
		IconURL:      fmt.Sprintf("/plugins/%s/profile.png", manifest.Id),
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

func (p *Plugin) getBotID() string {
	return p.BotUserID
}

func (p *Plugin) executeCommandHelp(args *model.CommandArgs) *model.CommandResponse {
	return p.responsef(args, getHelp(helpCommandText))
}

func responsef(format string, args ...interface{}) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf(format, args...),
		Type:         model.POST_DEFAULT,
	}
}

// ExecuteCommand executes a command that has been previously registered via the RegisterCommand API.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	if len(strings.Fields(args.Command)) < 2 {
		return p.executeCommandHelp(args), nil
	}

	split := strings.Fields(args.Command)
	action := split[1]

	switch action {
	case "add":
		return p.executeCommandAdd(args), nil
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

// executeCommandAdd adds a bookmark to the store
func (p *Plugin) executeCommandAdd(args *model.CommandArgs) *model.CommandResponse {
	subCommand := strings.Fields(args.Command)
	subCommand = subCommand[2:]

	if len(subCommand) < 1 {
		return p.responsef(args, "Missing sub-command. You can try %v", getHelp(addCommandText))
	}
	postID := p.getPostIDFromLink(subCommand[0])

	// verify postID exists
	post, appErr := p.API.GetPost(postID)
	if appErr != nil {
		return p.responsef(args, "PostID `%s` is not a valid postID", postID)
	}

	var bookmark Bookmark
	bookmark.PostID = postID

	// If no title provided, use the first X characters of the post message
	if len(subCommand) < 2 {
		numChars := math.Min(float64(len(post.Message)), MaxTitleCharacters)
		bookmark.Title = post.Message[0:int(numChars)]
	}

	p.addBookmark(args.UserId, &bookmark)

	return p.responsef(args, "Added bookmark: %+v", bookmark)
}

// executeCommandView shows all bookmarks in an ephemeral post
func (p *Plugin) executeCommandView(args *model.CommandArgs) *model.CommandResponse {
	bookmarks, err := p.getBookmarks(args.UserId)
	if err != nil {
		return p.responsef(args, "Unable to retrieve bookmarks for user %s", args.UserId)
	}

	if bookmarks == nil {
		return p.responsef(args, "You do not have any saved bookmarks")
	}

	team, appErr := p.API.GetTeam(args.TeamId)
	if appErr != nil {
		return p.responsef(args, "Unable to get team")
	}

	text := "#### Bookmarks List\n"
	for _, bmark := range bookmarks.ByID {
		text = text + p.bmarkBullet(bmark, team)
	}

	return p.responsef(args, text)
}

// executeCommandRemove removes a given bookmark from the store
func (p *Plugin) executeCommandRemove(args *model.CommandArgs) *model.CommandResponse {
	subCommand := strings.Fields(args.Command)

	if len(subCommand) < 3 {
		return p.responsef(args, "Missing sub-command. You can try %v", getHelp(removeCommandText))
	}

	bookmarkID := p.getPostIDFromLink(subCommand[2])

	err := p.deleteBookmark(args.UserId, bookmarkID)
	if err != nil {
		return p.responsef(args, err.Error())
	}

	return p.responsef(args, fmt.Sprintf("Removed bookmark ID: %s", bookmarkID))
}

func (p *Plugin) bmarkBullet(bmark *Bookmark, team *model.Team) string {
	return fmt.Sprintf("* %s - [%s](%s)\n", bmark.PostID, bmark.Title, p.getPermaLink(bmark.PostID, team.Name))
}

func (p *Plugin) getPermaLink(postID string, currentTeam string) string {
	return fmt.Sprintf("%v/%v/pl/%v", p.GetSiteURL(), currentTeam, postID)
}

func (p *Plugin) getPostIDFromLink(s string) string {
	r := regexp.MustCompile(`http:.*\/\w+\/\w+\/(\w+)`)
	if len(r.FindStringSubmatch(s)) == 2 {
		return r.FindStringSubmatch(s)[1]
	}
	return s
}
