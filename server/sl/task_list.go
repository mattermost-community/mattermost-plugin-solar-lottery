// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

func (sl *sl) ListTasks(rotation *Rotation, taskStatus TaskState) ([]string, error) {
	return []string{"<><> TODO"}, nil
}
