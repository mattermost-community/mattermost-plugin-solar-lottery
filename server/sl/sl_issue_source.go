// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type IssueService interface {
	DeleteIssueSource(rotationID, sourceName types.ID) error
	// PutIssueSource(rotationID types.ID, source *IssueSource) error
	MakeIssue(rotationID, sourceName types.ID) (*Task, error)
}

func (sl *sl) MakeIssue(rotationID types.ID, sourceName types.ID) (*Task, error) {
	var r *Rotation
	err := sl.Setup(
		pushLogger("MakeIssue", bot.LogContext{ctxSourceName: sourceName}),
		withRotation(rotationID, &r),
	)
	if err != nil {
		return nil, err
	}
	defer sl.popLogger()

	source, _ := r.IssueSource(sourceName)
	if source == nil {
		return nil, kvstore.ErrNotFound
	}

	t := source.NewTask()

	// Update Seq in the store
	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return nil, err
	}
	err = sl.Store.Entity(KeyTask).Store(t.TaskID, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (sl *sl) DeleteIssueSource(rotationID types.ID, sourceName types.ID) error {
	var r *Rotation
	err := sl.Setup(
		pushLogger("DeleteIssueSource", bot.LogContext{ctxSourceName: sourceName}),
		withRotation(rotationID, &r),
	)
	if err != nil {
		return nil
	}
	defer sl.popLogger()

	source, _ := r.IssueSource(sourceName)
	if source == nil {
		return kvstore.ErrNotFound
	}
	return sl.updateRotation(r, func() error {
		r.IssueSources.Delete(sourceName)
		return nil
	})
}

// func (sl *sl) SetIssueSource(r *Rotation, source *IssueSource) error {
// 	return sl.updateRotation(r, func() error {
// 		r.IssueSources.Set(source)
// 		return nil
// 	})
// }

// func (sl *sl) UpdateIssueSourceMin(r *Rotation, need Need, delete bool) error {
// 	return sl.updateRotation(r, func(r *Rotation) error {
// 		r.
// }
