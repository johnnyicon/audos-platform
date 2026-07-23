---
date: 2026-07-14
area: app-build
status: open
filed: no
label: Scheduled hooks never actually fire
---

Two recurring hourly schedules were created successfully (`POST /schedules`, valid RRULE, response
`status: pending`) against a hook confirmed working when called directly. Neither ever fired.

> Confirmed by direct observation: one schedule sat 2h07m past its `nextRun` time, the other 20m past —
> both still `runCount: 0`, `lastError: null`. Also found: one-time hook scheduling doesn't exist at all
> (`frequency` is a required field, rejecting any one-time payload); one-time scheduling only exists for
> the separate email-scheduling endpoint. Practical consequence: there is currently no working
> background/async execution path on this platform — anything slow has to finish inside a single hook's
> 5–10 minute ceiling, or be triggered from outside Audos entirely.
