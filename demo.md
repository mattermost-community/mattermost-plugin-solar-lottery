# Contents
- [Intro](#intro)
- [Demo Setup](#demo-setup)
- [Demo](#demo)
    - Init
        - [cleanup and initialize: `/lotto demo clean-init`](#lotto-demo-clean-init)
    - R&D Ice Breaker rotation
        - [initialize Ice Breaker: `/lotto demo icebreaker-init`](#lotto-demo-icebreaker-init)
        - [run forecast for Ice Breaker: `/lotto demo icebreaker-forecast`](#lotto-demo-icebreaker-forecast)
        - [set up and run autopilot for Ice Breaker: `/lotto demo icebreaker-autopilot`](#lotto-demo-icebreaker-autopilot)
        - [run user demo in Ice Breaker: `/lotto demo icebreaker-user`](#lotto-demo-icebreaker-user)
    - Sustaining Engineering Team (SET) rotation
        - [initialize SET: `/lotto demo SET-init`](#lotto-demo-set-init)
        - [run forecast for SET: `/lotto demo SET-forecast`](#lotto-demo-set-forecast)
- [Next Steps and TODO](#next-steps-and-todo)

## Intro

- Solar Lottery is a team rotation scheduler, inspired by [PagerDuty OnCall](https://www.pagerduty.com/platform/on-call-management/), and its predecessor the early amazon.com pager tool.
- Name from a Philip K. Dick novel "[Solar Lottery](https://en.wikipedia.org/wiki/Solar_Lottery)".
- The main motivation to develop was to automate the Sustaining Engineering Team (SET) schedulng.
- Not a traditional queue, scheduling is based on probabilities, exponentially increasing since the last serve time.
- Features (basic):
    - Users have skills, rotations have needs, match and constrain.
    - Grace periods after serving shifts, apply within the rotation.
    - User "unavailable" events.
    - Complete manual control over shifts, or "Autopilot"

### Use cases

- Automate existing:
    - R&D ice breakers.
    - OKR status review.
    - SET.
- New rotations:
    - Monthly tech talks.
    - Blog posts.
- Integrations team:
    - Rotating scrum master (triage, run the meeting, daily checks).
    - Ops duty - plugin review, releases, config on community.
    - Community duty.

### Limitations 

- Command-line only, optimized for debugging, not usability.
- Slow performance - unoptimized RPC access, O(N**2) forecasting.
- No cron yet - simulated for the demo with `/lotto rotation autopilot --debug-run <date>`.
- Practically no tests, the code can use a little refactoring.

### Demo notes

- Setting "unavailabe" times - supported but not explicitly demoed.
    - Grace after serving a shift.
    - Can create custom "Unavailable" intervals.
- Logs tee-ed as DMs.
    - Find very useful for debugging.
    - Started as debugging (got annoyed grepping logs).
    - INFOs are useful as a "transparency" change log, DEBUGs for customer support.

## Demo Setup

#### Team ABC (full stack)
- @aaron.medina - lead(3), server(3), webapp(3)
- @aaron.peterson - server(1), webapp(2), mobile(1)
- @aaron.ward - webapp(2)
- @albert.torres - server(1)
- @alice.johnston - server(2), webapp(1)

#### Team DEF (perf)
- @deborah.freeman - lead(2), server(3), webapp(2)
- @diana.wells - sever(3)
- @douglas.daniels - server(4)
- @emily.meyer - server(3), webapp(1)
- @eugene.rodriguez - server(2), webapp(1)
- @frances.elliott - webapp(3), mobile(3)

#### Team GHIJ (full stack)
- @helen.hunter - lead(3),webapp(1),server(1)
- @janice.armstrong - server(2), webapp(1)
- @jeremy.williamson - webapp(2)
- @jerry.ramos - server(2), webapp(2)
- @johnny.hansen - webapp(1), server(1), mobile(1)
- @jonathan.watson - server(2)

#### Team KL (SRE)
- @karen.austin - lead(2), sre(3)
- @karen.martin - webapp(1), sre(1), build(3)
- @kathryn.mills - server(2), sre(3)
- @laura.wagner - sre(2)

#### Team MN (full stack)
- @margaret.morgan - lead(3), server(2)
- @mark.rodriguez webapp(1), server(1)
- @matthew.mendoza webapp(3), server(2)
- @mildred.barnes webapp(2), server(2)
- @nancy.roberts webapp(2), server(3), mobile(1)

#### Team R (mobile)
- @ralph.watson - lead(2), mobile(3), webapp(2)
- @raymond.austin - mobile(2), webapp(1)
- @raymond.fisher - mobile(2), webapp(1)
- @raymond.fox - mobile(2), server(1), webapp(1)

## Demo

### `/lotto demo clean-init`

https://github.com/levb/mattermost-plugin-solar-lottery/blob/33ab0b9268a229d3c1ff4385df070d8bdd5b2d6f/server/plugin/demo.go#L15-L44

### `/lotto demo icebreaker-init`

```sh
/lotto rotation add -r icebreaker --period w --start 2019-12-17 --grace 3 --size 2

/lotto rotation join -r icebreaker -s 2019-12-17 -u @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston,@deborah.freeman,@diana.wells,@douglas.daniels,@emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@johnny.hansen,@jonathan.watson,@karen.austin,@karen.martin,@kathryn.mills,@laura.wagner,@margaret.morgan,@mark.rodriguez,@matthew.mendoza,@mildred.barnes,@nancy.roberts,@ralph.watson,@raymond.austin,@raymond.fisher,@raymond.fox
/lotto rotation show -r icebreaker
```

### `/lotto demo icebreaker-forecast`

```sh
/lotto rotation guess -r icebreaker -s 0 -n 20
/lotto rotation forecast -r icebreaker -s 0 -n 20 --sample 100
```

### `/lotto demo icebreaker-autopilot`

```sh
/lotto rotation autopilot -r icebreaker --notify 3 --fill-before 5
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-11
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-12
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-13
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-17
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-21
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-25
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-29
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-30
```

### `/lotto demo SET-init`

1. 6 "teams", 5 people in rotation, 1 grace shift after each serve.
2. require exactly 1 lead.
3. require a minimum of 2 each server and webapp, any level.
4. require exactly 1 mobile, any level.
5. require exactly 1 from each FS team.
6. require no minimum, max 1 from each of the other teams.

```sh
/lotto rotation add  -r SET --period m --start 2020-01-16 --grace 1 --size 5

/lotto rotation need -r SET --skill webapp --level beginner --min 2
/lotto rotation need -r SET --skill server --level beginner --min 2
/lotto rotation need -r SET --skill mobile --level beginner --min 1
/lotto rotation need -r SET --skill lead --level beginner --min 1 --max 1

/lotto rotation need -r SET --skill ABC-FS1 --level beginner --min 1 --max 1
/lotto rotation need -r SET --skill GHIJ-FS2 --level beginner --min 1 --max 1
/lotto rotation need -r SET --skill MN-FS3 --level beginner --min 1 --max 1
/lotto rotation need -r SET --skill DEF-PERF --level beginner --min -1 --max 1
/lotto rotation need -r SET --skill KL-SRE --level beginner --min -1 --max 1
/lotto rotation need -r SET --skill R-MOBILE --level beginner --min -1 --max 1

/lotto rotation join -r SET -u @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston,@deborah.freeman,@diana.wells,@douglas.daniels,@emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@johnny.hansen,@jonathan.watson,@karen.austin,@karen.martin,@kathryn.mills,@laura.wagner,@margaret.morgan,@mark.rodriguez,@matthew.mendoza,@mildred.barnes,@nancy.roberts,@ralph.watson,@raymond.austin,@raymond.fisher,@raymond.fox

/lotto rotation show -r SET

```

### `/lotto demo SET-forecast`

```sh
/lotto rotation guess -r SET -s 0 -n 3
/lotto rotation forecast -r SET -s 0 -n 3 --sample 5
```

## Next Steps and TODO
- **Intake** the plugin - review, transfer
- **Missing MVP features**:
    - **`/lotto shift leave`**, a little non-trivial for started shifts.
    - **Cron**, HA-aware.
    - **Import/Export** - for backups, especially while in beta
- **UI**, a more usable command set - intuitive, less verbose.
- **HTTP API**.
- **Caching** the RPC/KV access within a single request.
- **Submit to marketplace**
- **Features**:
    - A "Rotation Channel".
        - Tee relevant logs.
        - Automatically update the channel header.
    - History, archiving user events.
    - Past statistics.
    - Rotation name-spacing/isolation and admin access control.
    - Alternate scheduling strategies.
        - traditional queue.
        - better fairness, use all history not just last served date.
    - **◈** Integration with the **Welcome Bot**, **Workflow**, **Todo**, **Autolink** plugins.
    - **◈◈** Integration with the **Calendar** plugins.
