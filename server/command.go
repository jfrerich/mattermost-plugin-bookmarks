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
* |/bookmarks add <post_id> <bookmark_title> --labels <label1,label2>| - add a bookmark by specifying a post_id (with optional title)
* |/bookmarks add <permalink> <bookmark_title> --labels <label1,label2>| - add a bookmark by specifying the post permalink (with optional title)
`
	labelCommandText = `
**/bookmarks label**
* |/bookmarks label <post_id> --labels <labels>| - add labels (comma-separated) to a bookmark
* |/bookmarks label add <labels> | - create a new label
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
* |/bookmarks remove <post_id1>,<post_id2>| - remove multiple bookmarks by post_id, or permalink
`
	renameCommandText = `
**/bookmarks rename**
* |/bookmarks rename <label-old> <label-new>| - rename a label
`
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

func (p *Plugin) getBmarkTextOneLine(bmark *Bookmark, args *model.CommandArgs) (string, *model.AppError) {
	team, appErr := p.API.GetTeam(args.TeamId)
	if appErr != nil {
		return "", appErr
	}

	labelNames, _ := p.getLabelsForBookmark(args.UserId, bmark.PostID)
	// TODO: reconcile error types
	// if appErr != nil {
	// 	return "", appErr
	// }
	codeBlockedNames := p.getPrintableLabels(labelNames)

	titleFromPostLabel := ""
	title := bmark.Title
	if !bmark.hasUserTitle(bmark) {
		titleFromPostLabel = "`TitleFromPost` "
		title, appErr = p.getTitleFromPost(bmark)
		if appErr != nil {
			return "", appErr
		}
	}

	text := fmt.Sprintf("%s%s %s%s\n", p.getIconLink(bmark, team), codeBlockedNames, titleFromPostLabel, title)
	return text, nil
}

func (p *Plugin) getPrintableLabels(names []string) string {
	labels := ""
	for _, name := range names {
		labels = labels + fmt.Sprintf(" `%s`", name)
	}
	return labels
}

func (p *Plugin) getBmarkTextDetailed(bmark *Bookmark, args *model.CommandArgs) (string, *model.AppError) {
	team, appErr := p.API.GetTeam(args.TeamId)
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

	labelNames, _ := p.getLabelsForBookmark(args.UserId, bmark.PostID)
	// TODO: reconcile error types
	// if appErr != nil {
	// 	return "", appErr
	// }
	codeBlockedNames := p.getPrintableLabels(labelNames)
	post, appErr := p.API.GetPost(bmark.PostID)
	if appErr != nil {
		return "", appErr
	}

	iconLink := p.getIconLink(bmark, team)

	// team := post.
	text := fmt.Sprintf("%s\n#### Bookmark Title %s\n", codeBlockedNames, iconLink)
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
