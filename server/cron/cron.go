// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package cron

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/api"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
)

type cron struct {
	api    api.API
	config *config.Config
}
