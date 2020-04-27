# Mattermost Solar Lottery Plugin (work in progress, PRE-ALPHA)

[![CircleCI](https://circleci.com/gh/mattermost/mattermost-plugin-solar-lottery.svg?style=shield)](https://circleci.com/gh/mattermost/mattermost-plugin-solar-lottery)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattermost/mattermost-plugin-solar-lottery)](https://goreportcard.com/report/github.com/mattermost/mattermost-plugin-solar-lottery)
[![Code Coverage](https://img.shields.io/codecov/c/github/mattermost/mattermost-plugin-solar-lottery/master.svg)](https://codecov.io/gh/mattermost/mattermost-plugin-solar-lottery)

**Maintainer:** [@levb](https://github.com/levb)
**Co-Maintainer:** [@iomodo](https://github.com/iomodo)

A [Mattermost](https://mattermost.com) plugin somewhat similar to pager duty, allows to have rotations with magic "solar lottery" scheduling, or overrides.

## About

- Solar Lottery is a team rotation scheduler, inspired by [PagerDuty OnCall](https://www.pagerduty.com/platform/on-call-management/), and its predecessor the early amazon.com pager tool.
- Name from a Philip K. Dick novel "[Solar Lottery](https://en.wikipedia.org/wiki/Solar_Lottery)".
- The main motivation to develop was to automate the Sustaining Engineering Team (SET) schedulng.
- Not a traditional queue, scheduling is based on probabilities, exponentially increasing since the last serve time.
- Features (basic):
  - Users have skills, rotations have needs, match and constrain.
  - Grace periods after serving shifts, apply within the rotation.
  - User "unavailable" events.
  - Complete manual control over shifts, or "Autopilot"

## Install

1. Go the releases page and download the latest release.
2. On your Mattermost instance, go to System Console -> Plugin Management and upload it.
3. Configure plugin settings as desired.
4. Start using the plugin!

## Commands
  ### `/lotto rotation`
  Tools to manage rotations. 
  Usage: `/lotto <subcommand> <rotation-ID> [--flags]`.
  Subcommands: **archive** - **list** - **new** - **show** - **set** (autopilot|fill|limit|require|task)
  
  #### `/lotto rotation new`
  Creates a new rotation. Certain parameters can be specified only at creation
  time and may not be changed later.

  Flags:
    - `--beginning=datetime` - Beginning of time for shifts. Default: now.
    - `--fill-type=solar-lottery` - Task auto-assign type: only `solar-lottery` is currently supported,
    - `--fuzz int` - increase randomness of task assignment. Works by increasing the user weight doubling time by this many periods. Setting it above 3 will essentially make task assignemts random. Default: 0.
    - `--period=(weekly|biweekly|monthly)` - Recurrence period. For shifts, it is directly relevant; for tasks it affects how the user weights are calculated (shorter period leads to stricter rotation rules, much like lower fuzz). Default: `weekly`.
    - `--task-type=(shift|ticket)` - Currently, a rotation can only have _shifts_, i.e. recurring tasks, or _tickets_ that are submitted from an external source. Default: `shift`.
  
  #### `/lotto rotation archive`
  Archives a rotation.

  #### `/lotto rotation list`
  Lists active rotations.

  #### `/lotto rotation show`
  Shows rotation details.

  #### `/lotto rotation set autopilot`
  Changes rotation's autopilot settings.

  Flags:
      - `--off` - turns autopilot off for the rotation
      - [x] --create --create-prior[=28d]
      - [x] --schedule --schedule-prior[=7d]
      - [x] --start-finish
      - [x] --notify-start-prior[=3d]
      - [x] --notify-finish-prior[=3d]
      - [x] --run=time
    - [x] set fill
      - [ ] --beginning
      - [ ] --period
      - [ ] --seed
      - [ ] --fuzz
    - [x] set limit --skill <s-l> (--count | --clear)
    - [x] set require --skill <s-l> (--count | --clear)
    - [x] set task
      - [ ] --type=(shift|ticket)
      - [ ] --duration
      - [ ] --grace    
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


## By example

### Set u



## Usage

1. Start by creating a rotation with `/lotto rotation new` command
2. 
## Commands 
## Demo

