// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/pkg/errors"
)

type Issues interface {
	DeleteIssueSource(r *Rotation, sourceName string) error
	PutIssueSource(r *Rotation, source *IssueSource) error
	MakeIssue(r *Rotation, sourceName string) (*Task, error)
}

func (sl *sl) MakeIssue(r *Rotation, sourceName string) (*Task, error) {
	source, _ := r.IssueSource(sourceName)
	if source == nil {
		return nil, errors.New("Not found: " + sourceName)
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

func (sl *sl) DeleteIssueSource(r *Rotation, sourceName string) error {
	source, i := r.IssueSource(sourceName)
	if source == nil {
		return kvstore.ErrNotFound
	}
	return sl.UpdateRotation(r, func(r *Rotation) error {
		updated := r.IssueSources[:i]
		if i+1 < len(r.IssueSources) {
			updated = append(updated, r.IssueSources[i+1:]...)
		}
		return nil
	})
}

func (sl *sl) PutIssueSource(r *Rotation, source *IssueSource) error {
	found, i := r.IssueSource(source.Name)
	return sl.UpdateRotation(r, func(r *Rotation) error {
		if found == nil {
			r.IssueSources = append(r.IssueSources, source)
		} else {
			r.IssueSources[i] = source
		}
		return nil
	})
}
