// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Issues interface {
	DeleteIssueSource(r *Rotation, sourceName types.ID) error
	PutIssueSource(r *Rotation, source *IssueSource) error
	MakeIssue(r *Rotation, sourceName types.ID) (*Task, error)
}

func (sl *sl) MakeIssue(r *Rotation, sourceName types.ID) (*Task, error) {
	source, _ := r.IssueSource(sourceName)
	if source == nil {
		return nil, kvstore.ErrNotFound
	}

	t := source.NewTask()

	// Update Seq in the store
	err := sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
	if err != nil {
		return nil, err
	}

	err = sl.Store.Entity(KeyTask).Store(t.TaskID, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (sl *sl) DeleteIssueSource(r *Rotation, sourceName types.ID) error {
	source, _ := r.IssueSource(sourceName)
	if source == nil {
		return kvstore.ErrNotFound
	}
	return sl.updateExpandedRotation(r, func(r *Rotation) error {
		r.IssueSources.Delete(sourceName)
		return nil
	})
}

func (sl *sl) PutIssueSource(r *Rotation, source *IssueSource) error {
	return sl.updateExpandedRotation(r, func(r *Rotation) error {
		r.IssueSources.Set(source)
		return nil
	})
}
