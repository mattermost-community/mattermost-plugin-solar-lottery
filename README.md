# Mattermost Solar Lottery Plugin (work in progress, PRE-ALPHA)

[![CircleCI](https://circleci.com/gh/mattermost/mattermost-plugin-solar-lottery.svg?style=shield)](https://circleci.com/gh/mattermost/mattermost-plugin-solar-lottery)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattermost/mattermost-plugin-solar-lottery)](https://goreportcard.com/report/github.com/mattermost/mattermost-plugin-solar-lottery)
[![Code Coverage](https://img.shields.io/codecov/c/github/mattermost/mattermost-plugin-solar-lottery/master.svg)](https://codecov.io/gh/mattermost/mattermost-plugin-solar-lottery)

**Maintainer:** [@levb](https://github.com/levb)
**Co-Maintainer:** [@iomodo](https://github.com/iomodo)

A [Mattermost](https://mattermost.com) plugin somewhat similar to pager duty, allows to have rotations with magic "solar lottery" scheduling, or overrides.

## About

- Solar Lottery is a team rotation scheduler, inspired by [PagerDuty
  OnCall](https://www.pagerduty.com/platform/on-call-management/), and its
  predecessor the early amazon.com pager tool.
- Name from a Philip K. Dick novel "[Solar
  Lottery](https://en.wikipedia.org/wiki/Solar_Lottery)".
- The main motivation to develop was to automate the Sustaining Engineering Team
  (SET) schedulng.
- Not a traditional queue, scheduling is based on probabilities, exponentially
  increasing since the last serve time.
- Features (basic):
  - Users have skills, rotations have needs, match and constrain.
  - Grace periods after serving shifts, apply within the rotation.
  - User "unavailable" events.
  - Complete manual control over shifts, or "Autopilot"

## Install

1. Go the releases page and download the latest release.
2. On your Mattermost instance, go to System Console -> Plugin Management and
   upload it.
3. Configure plugin settings as desired.
4. Start using the plugin!

## Examples 

- ["Ice Breaker"](./server/command/use_case_ice_breaker_test.go) - a simple
rotation to select 2 individuals for a weeklt meeting's "ice breaker" 5 minute
intro.
- [SET](./server/command/use_case_set_test.go) - a monthly shift rotation with
  several skill requirements and constraints.
- [Autopilot](./server/command/rotation_autopilot_test.go) - an illustration of
  what happens when running autopilot.

## Commands

### `/lotto rotation`

Tools to manage rotations. 

Usage: `/lotto <subcommand> <rotation-ID> [--flags]`.

Subcommands: [archive](#lotto-rotation-archive) - [list](#lotto-rotation-list) - [new](#lotto-rotation-new) - [show](#lotto-rotation-show) - [set autopilot](#lotto-rotation-set-autopilot) | [set fill](#lotto-rotation-set-fill) | [set limit](#lotto-rotation-set-limit) | [set require](#lotto-rotation-set-require) | [set task](#lotto-rotation-set-task)

#### `/lotto rotation new`

Creates a new rotation. Certain parameters can be specified only at creation
time and may not be changed later.

Flags:

- `--beginning=datetime` - Beginning of time for shifts. Default: now.
- `--fill-type=solar-lottery` - Task auto-assign type: only `solar-lottery` is
  currently supported.
- `--fuzz int` - increase randomness of task assignment. Works by increasing the
  user weight doubling time by this many periods. Setting it above 3 will
  essentially make task assignemts random. Default: 0.
- `--period=(weekly|biweekly|monthly)` - Recurrence period. For shifts, it is
  directly relevant; for tasks it affects how the user weights are calculated
  (shorter period leads to stricter rotation rules, much like lower fuzz).
  Default: `weekly`.
- `--task-type=(shift|ticket)` - Currently, a rotation can only have _shifts_,
  i.e. recurring tasks, or _tickets_ that are submitted from an external source.
  Default: `shift`.

#### `/lotto rotation archive`

Archives a rotation.

#### `/lotto rotation list`

Lists active rotations.

#### `/lotto rotation show`

Shows rotation details.

#### `/lotto rotation set autopilot`

Changes rotation's autopilot settings.

Flags:

- `--off` - turns autopilot off for the rotation.
- `--create` - automatically create new shifts.
- `--create-prior=duration` - create new shifts this far ahead of their scheduled starts.
- `--schedule` - automatically schedule pending tasks.
- `--schedule-prior=duration` - schedule pending tasks this far ahead of their scheduled starts.
- `--start-finish` - automatically start and finish shifts that are due.
- `--remind-finish` - remind task users ahead ahead of its finish.
- `--remind-finish-prior` - remind this far ahead of the task finish.
- `--remind-start` - remind task users ahead of the start of a task.
- `--remind-start-prior` - remind this far ahead of the task start.

#### `/lotto rotation set fill`

Changes rotation's settings for filling (assigning users to) tasks.

Flags:

- `--fuzz` - adding fuzz slows down the exponential growth of idle users'
  weights, by adding this many rotation periods to the doubling time.

#### `/lotto rotation set limit`

Changes rotation's constraints (limits). A limit is like, "no more than 2 people
with knowledge of netops, intermediate plus. Use `any` to indicate any
skill/level.

Flags:

- `--skill=skill-level[,...]` - specifies the skill and the minimum level that the limit applies to. _skill_ can be any known skill, _level_ is a number 1-4, or omit the _-level_ to indicate that any level for the skill should count (same as setting to 1).
- `--count=number` - specifies the limit for the skill.
- `--clear` - clears the limit for the skill.

#### `/lotto rotation set require`

Changes rotation's requirements (needs). A requirement is like, "at least 2
people with knowledge of netops, intermediate plus. Use `any` to indicate any
skill/level.

Flags:

- `--skill=skill-level[,...]` - specifies the skill and the minimum level for the requirement. _skill_ can be any known skill, _level_ is a number 1-4, or omit the _-level_ to indicate that any level for the skill should count (same as setting to 1).
- `--count=number` - specifies how many users required for the skill.
- `--clear` - clears the requirement for the skill.

#### `/lotto rotation set task`

Changes rotation's defaults for new tasks.

Flags:
- `--duration` - sets the default duration for new tasks.
- `--grace` - sets the default grace period for new tasks.




((TODO))
- [ ] task
  - [ ] debug-delete
  - [ ] list --pending | --scheduled | --started | --finished
  - [x] assign
  - [x] close
  - [x] fill
  - [x] new shift
  - [x] new ticket
  - [x] schedule
  - [x] show ROT#id
  - [x] start
  - [x] unassign
- [x] user: manage users.
  - [x] disqualify [@user...] --skill 
  - [x] join ROT [@user...] --starting
  - [x] leave ROT [@user...]
  - [x] qualify [@user...] --skill 
  - [x] show [@user...]
  - [x] unavailable: [@user...] --start --finish [--clear] 
- [ ] autopilot [--now=datetime]
- [x] info: display this.
- [x] skill
  - [x] delete SKILL
  - [x] list
  - [x] new SKILL

