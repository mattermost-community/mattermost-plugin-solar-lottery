// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

type Expander interface {
	ExpandUserMap(UserMap) error
	ExpandUser(*User) error
	ExpandShift(*Shift) error
	ExpandRotation(*Rotation) error
}
