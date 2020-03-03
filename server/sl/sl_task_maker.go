// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

// func (sl *sl) MakeIssue(rotationID types.ID, sourceName types.ID) (*Task, error) {
// 	var r *Rotation
// 	err := sl.Setup(
// 		pushLogger("MakeIssue", bot.LogContext{ctxSourceName: sourceName}),
// 		withExpandedRotation(rotationID, &r),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer sl.popLogger()

// 	source, _ := r.IssueSource(sourceName)
// 	if source == nil {
// 		return nil, kvstore.ErrNotFound
// 	}

// 	t := source.NewTask()

// 	// Update Seq in the store
// 	err = sl.Store.Entity(KeyRotation).Store(r.RotationID, r)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = sl.Store.Entity(KeyTask).Store(t.TaskID, t)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return t, nil
// }

// func (sl *sl) UpdateIssueSourceMin(r *Rotation, need Need, delete bool) error {
// 	return sl.updateRotation(r, func(r *Rotation) error {
// 		r.
// }
