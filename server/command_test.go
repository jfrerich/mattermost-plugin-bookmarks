package main

import (
	"encoding/json"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecuteCommandView(t *testing.T) {

	bmarks1 := getTestBookmarks()
	jsonBmarks1, err := json.Marshal(bmarks1)
	assert.Nil(t, err)

	makePlugin := func(api *plugintest.API) *Plugin {
		p := &Plugin{}
		p.kvstore = NewStore(p)
		p.SetAPI(api)
		return p
	}

	u1 := &model.User{Id: model.NewId()}
	u2 := &model.User{Id: model.NewId()}

	api := makeAPIMock()
	siteURL := "https://myhost.com"
	p := makePlugin(api)

	api.On("GetTeam", mock.Anything).Return(&model.Team{Id: "teamID-1"}, nil)
	api.On("GetConfig", mock.Anything).Return(&model.Config{ServiceSettings: model.ServiceSettings{SiteURL: &siteURL}})

	api.On("getBookmarks", u1.Id).Return(bmarks1, nil)
	api.On("getBookmarks", u2.Id).Return(nil, nil)
	api.On("KVGet", getBookmarksKey(u1.Id)).Return(jsonBmarks1, nil)
	api.On("KVGet", getBookmarksKey(u2.Id)).Return(nil, nil)

	t.Run("no bookmarks", func(t *testing.T) {
		resp2 := p.executeCommandView(
			&model.CommandArgs{
				Command: "/bookmarks view",
				UserId:  u2.Id,
				TeamId:  "teamID-1",
			})
		assert.Contains(t, resp2.Text, "You do not have any saved bookmarks")
	})

	t.Run("3 bookmarks", func(t *testing.T) {
		resp := p.executeCommandView(
			&model.CommandArgs{
				Command: "/bookmarks view",
				UserId:  u1.Id,
				TeamId:  "teamID-1",
			})
		assert.Contains(t, resp.Text, "Title1")
		assert.Contains(t, resp.Text, "Title2")
		assert.Contains(t, resp.Text, "Title3")
	})

}

func makeAPIMock() *plugintest.API {
	api := &plugintest.API{}

	api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
	api.On("LogWarn", mock.Anything, mock.Anything, mock.Anything).Maybe()
	api.On("LogError", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

	return api
}
