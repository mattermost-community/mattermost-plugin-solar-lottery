// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type IssueSource struct {
	Name     string        `json:"name"`
	Seq      int           `json:"seq"`
	Requires Needs         `json:"requires,omitempty"`
	Limits   Needs         `json:"limits,omitempty"`
	Grace    time.Duration `json:"grace,omitempty"`
}

func NewIssueSource(name string) *IssueSource {
	return &IssueSource{
		Name: name,
		Seq:  1,
	}
}

func (is *IssueSource) NewTask() *Task {
	id := fmt.Sprintf("%s-%v", is.Name, is.Seq)
	is.Seq++

	t := NewTask()
	t.TaskID = types.ID(id)
	t.Requires = is.Requires
	t.Limits = is.Limits
	t.Grace = is.Grace

	return t
}

func (is IssueSource) MarkdownBullets() string {
	out := fmt.Sprintf("- %s\n", is.Markdown())
	out += fmt.Sprintf("  - Requires: **%s**\n", is.Requires.Markdown())
	out += fmt.Sprintf("  - Limits: **%v**\n", is.Limits.Markdown())
	if is.Grace != 0 {
		out += fmt.Sprintf("  - Grace: %v\n", is.Grace)
	}
	return out
}

func (is IssueSource) Markdown() string {
	return fmt.Sprintf("%s", is.Name)
}
