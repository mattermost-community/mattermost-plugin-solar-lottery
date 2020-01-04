## Intro

- Solar Lottery is a team rotation scheduler, inspired by PagerDuty (and Amazon.com pager tool before that).
- Name inspired by https://en.wikipedia.org/wiki/Solar_Lottery
- The main motivation to develop it was to automate the Sustaining Engineering Team (SET) schedulng.
- Not a traditional queue, scheduling is based on probabilities, exponentially increasing since last serve time.
- Features (basic):
    - Users have skills, rotations have needs, match and constrain.
    - Grace periods after serving shifts.
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

### Demo outline

- Scheduling the R&D "Ice Breakers".
- Scheduling SET.
- Setting "unavailabe" times.
    - Grace after serving a shift
    - Can create custom "Unavailable" intervals
- Logs tee-ed as DMs.
    - Started as debugging (got annoyed grepping logs).
    - INFOs are useful as a "transparency" change log.

## Demo Setup

#### Team ABC (full stack)
- @aaron.medina - lead(3), server(3), webapp(3)
- @aaron.peterson - server(1), webapp(2)
- @aaron.ward - webapp(2)
- @albert.torres - server(1)
- @alice.johnston - server(2), webapp(1)

#### Team DEF (perf)
- @deborah.freeman - lead(2), server(3), webapp(2)
- @diana.wells - sever(3)
- @douglas.daniels - server(4)
- @emily.meyer - server(3), webapp(1)
- @eugene.rodriguez - server(2), webapp(1)
- @frances.elliott - webapp(3)

#### Team GHIJ (full stack)
- @helen.hunter - lead(3),webapp(1),server(1)
- @janice.armstrong - server(2), webapp(1)
- @jeremy.williamson - webapp(2)
- @jerry.ramos - server(2), webapp(2)
- @johnny.hansen - webapp(1), server(1)
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
- @nancy.roberts webapp(2), server(3)

## Demo

### Setup
(done beforehand)

```sh
/lotto debug-clean

/lotto user qualify -k ABC -l intermediate -u @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston
/lotto user qualify -k DEF -l intermediate -u @deborah.freeman,@diana.wells,@douglas.daniels,@emily.meyer,@eugene.rodriguez,@frances.elliott
/lotto user qualify -k GHIJ -l intermediate -u @helen.hunter,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@johnny.hansen,@jonathan.watson
/lotto user qualify -k KL -l intermediate -u @karen.austin,@karen.martin,@kathryn.mills,@laura.wagner
/lotto user qualify -k MN -l intermediate -u @margaret.morgan,@mark.rodriguez,@matthew.mendoza,@mildred.barnes,@nancy.roberts

/lotto user qualify -k lead -l intermediate -u @deborah.freeman,@karen.austin
/lotto user qualify -k lead -l advanced -u @aaron.medina,@helen.hunter,@margaret.morgan

/lotto user qualify -k server -l beginner -u @aaron.peterson,@albert.torres,@helen.hunter,@johnny.hansen,@mark.rodriguez
/lotto user qualify -k server -l intermediate -u @alice.johnston,@eugene.rodriguez,@janice.armstrong,@jerry.ramos,@jonathan.watson,@kathryn.mills,@margaret.morgan,@matthew.mendoza,@mildred.barnes
/lotto user qualify -k server -l advanced -u @aaron.medina,@deborah.freeman,@diana.wells,@emily.meyer,@nancy.roberts
/lotto user qualify -k server -l expert -u @douglas.daniels

/lotto user qualify -k webapp -l beginner -u @emily.meyer,@eugene.rodriguez,@helen.hunter,@janice.armstrong,@johnny.hansen,@karen.martin,@mark.rodriguez,@mildred.barnes
/lotto user qualify -k webapp -l intermediate -u @aaron.peterson,@aaron.ward,@deborah.freeman,@jeremy.williamson,@jerry.ramos,@nancy.roberts
/lotto user qualify -k webapp -l advanced -u @aaron.medina,@frances.elliott,@matthew.mendoza

/lotto user qualify -k sre -l beginner -u @karen.martin
/lotto user qualify -k sre -l intermediate -u @laura.wagner
/lotto user qualify -k sre -l advanced -u @karen.austin,@kathryn.mills

/lotto user qualify -k build -l advanced -u @karen.martin

/lotto rotation add -r icebreaker --period w --start 2019-12-17 --grace 3 --size 2

/lotto rotation join -r icebreaker -s 2019-12-17 -u @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston,@deborah.freeman,@diana.wells,@douglas.daniels,@emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@johnny.hansen,@jonathan.watson,@karen.austin,@karen.martin,@kathryn.mills,@laura.wagner,@margaret.morgan,@mark.rodriguez,@matthew.mendoza,@mildred.barnes,@nancy.roberts

```

### 1. Ice-breaker
```sh
/lotto rotation guess -r icebreaker -s 0 -n 20 --autofill
/lotto rotation forecast -r icebreaker -s 0 -n 20 --sample 100
/lotto rotation autopilot -r icebreaker --notify 3 --fill-before 5
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-11
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-17
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-22
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-30
```

### 2. Sustaining Engineering Team (SET)

```sh
/lotto rotation add  -r SET --period m --start 2020-01-16 --grace 2 --size 5
/lotto rotation need -r SET --skill webapp --level beginner --min 1
/lotto rotation need -r SET --skill webapp --level intermediate --min 1 
/lotto rotation need -r SET --skill server --level beginner --min 1
/lotto rotation need -r SET --skill server --level intermediate --min 1
/lotto rotation need -r SET --skill lead --level beginner --min 1 --max 1
/lotto rotation need -r SET --skill ABC --level beginner --min 1 --max 1
/lotto rotation need -r SET --skill DEF --level beginner --min 1 --max 1
/lotto rotation need -r SET --skill GHIJ --level beginner --min 1 --max 1
/lotto rotation need -r SET --skill MN --level beginner --min 1 --max 1
# No KL - it is not required to contribute, only opportunistically

/lotto rotation join -r SET -u @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston,@deborah.freeman,@diana.wells,@douglas.daniels,@emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@johnny.hansen,@jonathan.watson,@karen.austin,@karen.martin,@kathryn.mills,@laura.wagner,@margaret.morgan,@mark.rodriguez,@matthew.mendoza,@mildred.barnes,@nancy.roberts

/lotto rotation show -r SET

/lotto rotation guess -r SET -s 0 -n 3 --autofill
/lotto rotation forecast -r SET -s 0 -n 10 --sample 5
```

### TODO
- **Cron**, HA-aware.
- **UI**, a more usable command set - intuitive, less verbose.
- **HTTP API**.
- **Caching** the RPC/KV access within a single request.
- **Ambitious** Integration with the **Welcome Bot**, **Workflow**, **Todo**, **Autolink** plugins.
- **Ambitious** Integration with the **Calendar** plugins.
- **Features**:
    - "Rotation Channel"
        - Tee relevant logs
        - Automatically update the channel header
    - History, archiving user events.
    - Past statistics.
    - Rotation name-spacing/isolation and admin access control.
    - Alternate scheduling strategies - traditional queue.