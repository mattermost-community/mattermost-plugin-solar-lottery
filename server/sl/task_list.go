// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import "github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"

func (sl *sl) ListTasks(rotation *Rotation, taskStatus types.ID) ([]string, error) {
	return []string{"<><> TODO"}, nil
}
