// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Rotation struct {
	PluginVersion string
	RotationID    types.ID
	AutofillType  string
	IsArchived    bool
	IssueSources  *types.Index // of *IssueSource
	// ShiftSource *ShiftSource
	Pending           []*Task
	InProgress        []*Task
	MattermostUserIDs *types.IDIndex `json:",omitempty"`

	users Users
}

func NewRotation() *Rotation {
	r := &Rotation{}
	r.init()
	return r
}

func (r *Rotation) init() {
	if r.MattermostUserIDs == nil {
		r.MattermostUserIDs = types.NewIDIndex()
	}
	if r.users.IsEmpty() {
		r.users = NewUsers()
	}
}

func (rotation *Rotation) WithMattermostUserIDs(pool Users) *Rotation {
	newRotation := *rotation
	newRotation.MattermostUserIDs = types.NewIDIndex()
	for _, id := range pool.IDs() {
		newRotation.MattermostUserIDs.Set(id)
	}
	if pool.IsEmpty() {
		pool = NewUsers()
	}
	newRotation.users = pool
	return &newRotation
}

func (r *Rotation) String() string {
	return r.Name()
}

func (r *Rotation) Name() string {
	return kvstore.NameFromID(r.RotationID)
}

func (r *Rotation) Markdown() string {
	return r.Name()
}

func (r *Rotation) MarkdownBullets() string {
	out := fmt.Sprintf("- **%s**\n", r.Name())
	out += fmt.Sprintf("  - ID: `%s`.\n", r.RotationID)
	out += fmt.Sprintf("  - Users (%v): %s.\n", r.MattermostUserIDs.Len(), r.users.MarkdownWithSkills())

	// if rotation.Autopilot.On {
	// 	out += fmt.Sprintf("  - Autopilot: **on**\n")
	// 	out += fmt.Sprintf("    - Auto-start: **%v**\n", rotation.Autopilot.StartFinish)
	// 	out += fmt.Sprintf("    - Auto-fill: **%v**, %v days prior to start\n", rotation.Autopilot.Fill, rotation.Autopilot.FillPrior)
	// 	out += fmt.Sprintf("    - Notify users in advance: **%v**, %v days prior to transition\n", rotation.Autopilot.Notify, rotation.Autopilot.NotifyPrior)
	// } else {
	// 	out += fmt.Sprintf("  - Autopilot: **off**\n")
	// }

	return out
}

func (r *Rotation) FindUsers(mattermostUserIDs *types.IDIndex) []*User {
	uu := []*User{}
	for _, id := range r.MattermostUserIDs.IDs() {
		uu = append(uu, r.users.Get(id))
	}
	return uu
}

func (r *Rotation) IssueSource(sourceName types.ID) (*IssueSource, int) {
	for i, id := range r.IssueSources.IDs() {
		if id == sourceName {
			return r.IssueSources.Get(id).(*IssueSource), i
		}
	}
	return nil, -1
}
