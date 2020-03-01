// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type IssueSource struct {
	Name  types.ID      `json:"name"`
	Seq   int           `json:"seq"`
	Min   Needs         `json:"min,omitempty"`
	Max   Needs         `json:"max,omitempty"`
	Grace time.Duration `json:"grace,omitempty"`
}

func NewIssueSource(name types.ID) *IssueSource {
	return &IssueSource{
		Name: name,
		Seq:  1,
		Min:  NewNeeds(),
		Max:  NewNeeds(),
	}
}

func (is *IssueSource) GetID() types.ID { return is.Name }

func (is *IssueSource) NewTask() *Task {
	id := fmt.Sprintf("%s-%v", is.Name, is.Seq)
	is.Seq++

	t := NewTask()
	t.TaskID = types.ID(id)
	t.Min = is.Min
	t.Max = is.Max
	t.Grace = is.Grace

	return t
}

func (is IssueSource) MarkdownBullets() string {
	out := fmt.Sprintf("- %s\n", is.Markdown())
	out += fmt.Sprintf("  - Requires: **%s**\n", is.Min.Markdown())
	out += fmt.Sprintf("  - Limits: **%v**\n", is.Max.Markdown())
	if is.Grace != 0 {
		out += fmt.Sprintf("  - Grace: %v\n", is.Grace)
	}
	return out
}

func (is IssueSource) Markdown() string {
	return fmt.Sprintf("%s", is.Name)
}

type issueSources []*IssueSource
