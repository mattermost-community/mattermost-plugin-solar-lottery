// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/kvstore"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Rotation struct {
	PluginVersion  string
	RotationID     types.ID
	IsArchived     bool
	TaskFillerType types.ID
	TaskMaker      *TaskMaker
	Starts         types.Time

	MattermostUserIDs *types.IDSet `json:",omitempty"`
	TaskIDs           *types.IDSet `json:",omitempty"`

	loaded bool
	Users  *Users `json:"-"`
	tasks  *Tasks
}

func NewRotation() *Rotation {
	r := &Rotation{
		Starts: types.NewTime(),
	}
	r.init()
	return r
}

func (r *Rotation) init() {
	if r.MattermostUserIDs == nil {
		r.MattermostUserIDs = types.NewIDSet()
	}
	if r.TaskIDs == nil {
		r.TaskIDs = types.NewIDSet()
	}
	if r.TaskMaker == nil {
		r.TaskMaker = NewTaskMaker()
	}
}

func (rotation *Rotation) WithMattermostUserIDs(pool *Users) *Rotation {
	newRotation := *rotation
	newRotation.MattermostUserIDs = types.NewIDSet()
	for _, id := range pool.IDs() {
		newRotation.MattermostUserIDs.Set(id)
	}
	if pool.IsEmpty() {
		pool = NewUsers()
	}
	newRotation.Users = pool
	return &newRotation
}

func (r *Rotation) String() string {
	return r.Name()
}

func (r *Rotation) Name() string {
	return kvstore.NameFromID(r.RotationID)
}

func (r *Rotation) Markdown() md.MD {
	return md.MD(r.Name())
}

func (r *Rotation) MarkdownBullets() md.MD {
	out := md.Markdownf("- **%s**\n", r.Name())
	out += md.Markdownf("  - ID: `%s`.\n", r.RotationID)
	if r.Users != nil {
		out += md.Markdownf("  - Users (%v): %s.\n", r.MattermostUserIDs.Len(), r.Users.MarkdownWithSkills())
	} else {
		out += md.Markdownf("  - Users (%v): %s.\n", r.MattermostUserIDs.Len(), r.MattermostUserIDs.IDs())
	}
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

func (r *Rotation) FindUsers(mattermostUserIDs *types.IDSet) []*User {
	uu := []*User{}
	for _, id := range r.MattermostUserIDs.IDs() {
		uu = append(uu, r.Users.Get(id))
	}
	return uu
}
