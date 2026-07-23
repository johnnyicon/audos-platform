# SOP: dispatch, then poll — the standard pattern for every Audos job

*Written 2026-07-13. This is the operating pattern that made today's work reliable instead of a series*
*of blocked, timed-out, or half-finished interactions. Use it for every job dispatched to Otto/Cursor,*
*no exceptions.*

## The problem this solves

Audos build/edit jobs (delegated to the Cursor Background Agent backend — see
`21-cursor-backend-research.md`) routinely take **anywhere from under a minute to 25+ minutes**. A
synchronous "send request, wait for the response" call either times out, blocks everything else you
could be doing, or both. Meanwhile, checking too early just gets you "Running" over and over, burning
calls for no information.

## The pattern

1. **Dispatch in the background.** Fire the job-creation request (a message to Otto instructing it to
   create a draft and run it by its own specific task ID — never run-all) as a background process, not
   a blocking call. Read the immediate response only for the **task ID** and dispatch confirmation —
   it will not contain the result yet.
2. **Immediately start a poller in the background**, separate from the dispatch call. The poller loops:
   sleep, ask Otto for a compact status line on that specific task ID, check if the response indicates
   `Complete`/`Failed` (not `Queued`/`Running`), and stop looping once it does. Write each poll's raw
   response to a log file — don't parse and discard, so you can inspect exactly what came back later if
   something's ambiguous.
3. **Keep working while it polls.** The whole point is not blocking — dispatch the next independent job,
   update documentation, or do anything else that doesn't depend on this result while the poller runs in
   the background.
4. **When the poller reports completion, read the final report — then verify live anyway.** A "Complete"
   status is not evidence the change is correct; see `19-capabilities-and-limitations.md`'s "self-reports
   are unreliable" and "check the initial paint" sections. The poller tells you *when* to look, not
   *what's true*.

## Concrete sleep interval

**90 seconds between polls**, up to ~15–20 iterations (so roughly 20–30 minutes of patience before
declaring a poller stuck). Faster than 90s wastes calls on jobs that reliably take several minutes
minimum; slower risks sitting idle long after a fast job actually finished. This was tuned empirically
over many jobs today, not picked arbitrarily.

## Why this beats the alternatives

- **Synchronous wait-and-block**: ties up the whole session for the slowest possible job, every time,
  even when 80% of jobs finish in under 3 minutes.
- **Checking immediately, no poller**: near-certain "Running" response, wasted round trip, and a strong
  temptation to just guess at the outcome instead of actually waiting for it.
- **Polling too aggressively (e.g. every 10–15s)**: burns far more calls than necessary for no
  informational gain — job state doesn't change that fast.

## What to log per dispatch

At minimum: timestamp, one-line description of what was asked, the model pinned (see the model-selection
rule in `AGENTS.md`), the task ID once known, and the final result once the poller completes. The running
record is kept in exactly this shape at `doknow-kb/audos/ACTIVITY-LOG.md` (a different repo — deliberately
not a file in this one, see `24-where-new-findings-go.md`) — use the `audos-sdk-log-activity` skill rather
than editing it by hand.
