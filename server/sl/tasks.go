// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"strings"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type Tasks struct {
	*types.ValueSet // of *Task
}

func NewTasks(tt ...*Task) *Tasks {
	tasks := &Tasks{
		ValueSet: types.NewValueSet(&taskArray{}),
	}
	for _, t := range tt {
		tasks.Set(t)
	}
	return tasks
}

func (tasks Tasks) Get(id types.ID) *Task {
	return tasks.ValueSet.Get(id).(*Task)
}

func (tasks Tasks) Markdown() string {
	out := []string{}
	for _, t := range tasks.AsArray() {
		out = append(out, t.Markdown())
	}
	return strings.Join(out, ", ")
}

func (tasks Tasks) String() string {
	out := []string{}
	for _, t := range tasks.AsArray() {
		out = append(out, t.String())
	}
	return strings.Join(out, ", ")
}

// TestArray returns all tasks, sorted by MattermostUserID. It is used in
// testing, so it returns []Task rather than a []*Task to make it easier to
// compare with expected results.
func (tasks Tasks) TestArray() []Task {
	out := []Task{}
	for _, id := range tasks.TestIDs() {
		t := tasks.Get(types.ID(id))
		out = append(out, *t)
	}
	return out
}

func (tasks Tasks) AsArray() []*Task {
	a := taskArray{}
	tasks.ValueSet.AsArray(&a)
	return []*Task(a)
}

type taskArray []*Task

func (p taskArray) Len() int                   { return len(p) }
func (p taskArray) GetAt(n int) types.Value    { return p[n] }
func (p taskArray) SetAt(n int, v types.Value) { p[n] = v.(*Task) }

func (p taskArray) InstanceOf() types.ValueArray {
	inst := make(taskArray, 0)
	return &inst
}
func (p *taskArray) Ref() interface{} { return &p }
func (p *taskArray) Resize(n int) {
	*p = make(taskArray, n)
}
