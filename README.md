# Mattermost Solar Lottery Plugin

[![CircleCI](https://circleci.com/gh/mattermost/mattermost-plugin-solar-lottery.svg?style=shield)](https://circleci.com/gh/mattermost/mattermost-plugin-solar-lottery)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattermost/mattermost-plugin-solar-lottery)](https://goreportcard.com/report/github.com/mattermost/mattermost-plugin-solar-lottery)
[![Code Coverage](https://img.shields.io/codecov/c/github/mattermost/mattermost-plugin-solar-lottery/master.svg)](https://codecov.io/gh/mattermost/mattermost-plugin-solar-lottery)

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

## Commands / Usage / Demo

For a comprehensive list of commands and use cases please See [Demo](demo.md)
