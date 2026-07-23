---
date: 2026-07-16
area: app-build
status: inconclusive
label: "Scheduled-hooks retest: was round 1's \"never fires\" finding cadence-specific?"
---

**Hypothesis:** Round 1 found two hourly scheduled hooks that never fired. Was that a genuine platform
failure, or specific to that cadence/timing?

**Method:** Created a daily (rather than hourly) recurring schedule and checked its firing time and
`runCount`; asked Otto directly whether the earlier non-fire was a known issue.

**Result: contradicts round 1, root cause left genuinely open.** The daily schedule fired correctly,
within about 9 seconds of its `nextRun` time — the opposite of round 1's result on the same underlying
mechanism. Whether the difference is cadence-specific, already fixed platform-side between the two
tests, or scoped per-workspace was left explicitly unresolved: "genuinely open, not closed by either
test alone." Otto had no documented knowledge of the discrepancy either way.

See `BACKLOG.md #12`; `CHANGELOG.md`, scheduled-hooks section.
