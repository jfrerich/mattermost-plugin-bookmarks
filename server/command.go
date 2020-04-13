package main

import (
	"fmt"
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

	split := strings.Fields(args.Command)
	if len(split) < 2 {
		return p.executeCommandHelp(args), nil
	}

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
	_, appErr := p.API.GetPost(postID)
	if appErr != nil {
		return p.responsef(args, "PostID `%s` is not a valid postID", postID)
	}

	var bookmark Bookmark
	bookmark.PostID = postID

	// Only save title if user provides one.
	if len(subCommand) >= 2 {
		bookmark.Title = subCommand[1]
	}

	p.addBookmark(args.UserId, &bookmark)

	text, appErr := p.getBmarkTextOneLine(&bookmark, args.TeamId)
	if appErr != nil {
		return p.responsef(args, "Unable to get bookmarks list bookmark")
	}

	return p.responsef(args, "Added bookmark: %s", text)
}

// executeCommandView shows all bookmarks in an ephemeral post
func (p *Plugin) executeCommandView(args *model.CommandArgs) *model.CommandResponse {

	subCommand := strings.Fields(args.Command)

	// user requests to view an indiviual bookmark
	if len(subCommand) == 3 {
		postID := subCommand[2]
		postID = p.getPostIDFromLink(postID)
		bmark, err := p.getBookmark(args.UserId, postID)
		if err != nil {
			return p.responsef(args, "Unable to retrieve bookmark for user %s", args.UserId)
		}

		text, appErr := p.getBmarkTextDetailed(bmark, args.TeamId)
		if appErr != nil {
			return p.responsef(args, "Unable to retrieve bookmark for user %s", args.UserId)
		}
		return p.responsef(args, text)
	}

	bookmarks, err := p.getBookmarks(args.UserId)
	if err != nil {
		return p.responsef(args, "Unable to retrieve bookmarks for user %s", args.UserId)
	}

	if bookmarks == nil {
		return p.responsef(args, "You do not have any saved bookmarks")
	}

	text := "#### Bookmarks List\n"
	for _, bmark := range bookmarks.ByID {
		nextText, appErr := p.getBmarkTextOneLine(bmark, args.TeamId)
		if appErr != nil {
			return p.responsef(args, "Unable to get bookmarks list bookmark")
		}
		text = text + nextText
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

	bmark, err := p.deleteBookmark(args.UserId, bookmarkID)
	if err != nil {
		return p.responsef(args, err.Error())
	}

	text, appErr := p.getBmarkTextOneLine(bmark, args.TeamId)
	if appErr != nil {
		return p.responsef(args, "Unable to get bookmarks list bookmark")
	}
	return p.responsef(args, fmt.Sprintf("Removed bookmark: %s", text))
}

// getTitleFromPost returns a title generated from a Post.Message
func (p *Plugin) getTitleFromPost(bmark *Bookmark) (string, *model.AppError) {

	// TODO: set limit to number of character from post.Message
	// numChars := math.Min(float64(len(post.Message)), MaxTitleCharacters)
	// bookmark.Title = post.Message[0:int(numChars)]
	post, appErr := p.API.GetPost(bmark.PostID)
	if appErr != nil {
		return "", appErr
	}
	title := post.Message
	return title, nil
}

func (p *Plugin) getBmarkTextOneLine(bmark *Bookmark, teamID string) (string, *model.AppError) {
	team, appErr := p.API.GetTeam(teamID)
	if appErr != nil {
		return "", appErr
	}

	titleFromPostLabel := ""
	title := bmark.Title
	if !bmark.hasUserTitle(bmark) {
		titleFromPostLabel = "`TitleFromPost` "
		title, appErr = p.getTitleFromPost(bmark)
		if appErr != nil {
			return "", appErr
		}
	}

	text := fmt.Sprintf("%s %s%s\n", p.getIconLink(bmark, team), titleFromPostLabel, title)
	return text, nil
}

func (p *Plugin) getBmarkTextDetailed(bmark *Bookmark, teamID string) (string, *model.AppError) {
	team, appErr := p.API.GetTeam(teamID)
	if appErr != nil {
		return "", appErr
	}

	title, appErr := p.getTitleFromPost(bmark)
	if appErr != nil {
		return "", appErr
	}

	if bmark.hasUserTitle(bmark) {
		title = bmark.Title
	}

	post, appErr := p.API.GetPost(bmark.PostID)
	if appErr != nil {
		return "", appErr
	}

	iconLink := p.getIconLink(bmark, team)

	// team := post.
	text := fmt.Sprintf("#### Bookmark Title %s\n", iconLink)
	text = text + fmt.Sprintf("**%s**\n", title)
	text = text + "##### Post Message \n"
	text = text + post.Message

	return text, appErr

}

func (p *Plugin) getIconLink(bmark *Bookmark, team *model.Team) string {
	iconLink := fmt.Sprintf("[:link:](%s)", p.getPermaLink(bmark.PostID, team.Name))
	return iconLink
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
