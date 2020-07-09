package main

import (
	"sync"

	"github.com/gorilla/mux"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/command"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	// BotId of the created bot account.
	BotUserID string

	router *mux.Router
}

// OnActivate runs when the plugin activates and ensures the plugin is properly
// configured.
func (p *Plugin) OnActivate() error {
	bot := &model.Bot{
		Username:    "bookmarks",
		DisplayName: "Bookmarks",
		Description: "Created by the Bookmarks plugin.",
	}
	options := []plugin.EnsureBotOption{
		plugin.ProfileImagePath("assets/profile.png"),
	}

	p.initialiseAPI()

	botID, err := p.Helpers.EnsureBot(bot, options...)
	if err != nil {
		return errors.Wrap(err, "failed to ensure Bookmarks bot")
	}
	p.BotUserID = botID

	// return p.API.RegisterCommand(createBookmarksCommand())
	command.Register(p.API.RegisterCommand)
	return nil
}

// GetSiteURL returns the SiteURL from the config settings
func (p *Plugin) GetSiteURL() string {
	ptr := p.API.GetConfig().ServiceSettings.SiteURL
	if ptr == nil {
		return ""
	}
	return *ptr
}

// ExecuteCommand executes a command that has been previously registered via the RegisterCommand API.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	pluginapi := pluginapi.New(p.API)
	command := command.Command{
		Context:   c,
		Args:      args,
		ChannelID: args.ChannelId,
		API:       pluginapi,
	}

	out := command.Handle()
	// if err != nil {
	// 	p.API.LogError(err.Error())
	// 	return nil, model.NewAppError("bookmarks.ExecuteCommand", "Unable to execute command.", nil, err.Error(), http.StatusInternalServerError)
	// }

	// if out != "" {
	// }

	post := &model.Post{
		UserId:    p.GetBotID(),
		ChannelId: args.ChannelId,
		Message:   out,
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)

	response := &model.CommandResponse{}

	return response, nil
}
