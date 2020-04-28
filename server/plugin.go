package main

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
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

	botID, err := p.Helpers.EnsureBot(bot, options...)
	if err != nil {
		return errors.Wrap(err, "failed to ensure Bookmarks bot")
	}
	p.BotUserID = botID

	return p.API.RegisterCommand(getCommand())
}

// GetSiteURL returns the SiteURL from the config settings
func (p *Plugin) GetSiteURL() string {
	ptr := p.API.GetConfig().ServiceSettings.SiteURL
	if ptr == nil {
		return ""
	}
	return *ptr
}
