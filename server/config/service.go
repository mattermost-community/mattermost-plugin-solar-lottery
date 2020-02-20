package config

import (
	"errors"
	"net/url"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
)

type Store interface {
	SaveConfig(conf Mapper)
	GetConfig(ref interface{}) error
	GetMattermostConfig() *model.Config
	GetBotUserID() string
}

type Mapper interface {
	Map(onto map[string]interface{}) (result map[string]interface{})
}

type Service interface {
	Get() *Config
	Refresh() error
	Store(*StoredConfig)
}

type service struct {
	*BuildConfig

	lock   *sync.RWMutex
	config *Config
	store  Store
}

func NewService(build *BuildConfig, store Store) Service {
	return &service{
		lock:        &sync.RWMutex{},
		store:       store,
		BuildConfig: build,
	}
}

func (s *service) Get() *Config {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.config == nil {
		return &Config{}
	}
	return &(*s.config)
}

func (s *service) Refresh() error {
	stored := StoredConfig{}
	err := s.store.GetConfig(&stored)
	if err != nil {
		return err
	}

	mattermostSiteURL := s.store.GetMattermostConfig().ServiceSettings.SiteURL
	if mattermostSiteURL == nil {
		return errors.New("plugin requires Mattermost Site URL to be set")
	}
	mattermostURL, err := url.Parse(*mattermostSiteURL)
	if err != nil {
		return err
	}
	pluginURLPath := "/plugins/" + s.BuildConfig.PluginID
	pluginURL := strings.TrimRight(*mattermostSiteURL, "/") + pluginURLPath

	s.lock.Lock()
	defer s.lock.Unlock()

	newConfig := s.config
	if newConfig == nil {
		newConfig = &Config{}
	}
	newConfig.StoredConfig = &stored
	newConfig.MattermostSiteURL = *mattermostSiteURL
	newConfig.MattermostSiteHostname = mattermostURL.Hostname()
	newConfig.PluginURL = pluginURL
	newConfig.PluginURLPath = pluginURLPath
	newConfig.BotUserID = s.store.GetBotUserID()

	return nil
}

func (s *service) Store(newStored *StoredConfig) {
	s.lock.Lock()
	s.config.StoredConfig = newStored
	// no run-time fields to update
	s.lock.Unlock()

	s.store.SaveConfig(s.config.StoredConfig)
}
