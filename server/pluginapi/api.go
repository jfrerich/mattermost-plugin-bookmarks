package pluginapi

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type api struct {
	papi plugin.API
}

type API interface {
	GetPost(postID string) (*model.Post, error)
	GetConfig() *model.Config
	KVSet(key string, value []byte) error
	KVGet(key string) ([]byte, error)
}

func New(a plugin.API) API {
	return &api{
		papi: a,
	}
}

func (a *api) GetPost(postID string) (*model.Post, error) {
	p, appErr := a.papi.GetPost(postID)
	if appErr != nil {
		return nil, appErr
	}
	return p, nil
}

func (a *api) KVSet(key string, value []byte) error {
	appErr := a.papi.KVSet(key, value)
	if appErr != nil {
		return appErr
	}
	return nil
}

func (a *api) KVGet(key string) ([]byte, error) {
	value, appErr := a.papi.KVGet(key)
	if appErr != nil {
		return nil, appErr
	}
	return value, nil
}

func (a *api) GetConfig() *model.Config {
	return a.papi.GetConfig()
}
