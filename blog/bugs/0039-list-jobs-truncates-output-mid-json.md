---
date: 2026-08-05
area: app-build
status: open
filed: no
label: The jobs list truncates a job's output mid-JSON, so diagnostic results come back unreadable
---

A Cursor delegation job (#98454) was dispatched to run three catalogue `SELECT`s and log the raw JSON
result. The job completed successfully and did exactly that. Reading the result back through the jobs list
returned the output **cut off mid-JSON**:

```
=== EXT-CHECK RAW OUTPUT ===
{"installed":{"ok":true,"rows":[{"extname":"plpgsql"}]},
 "available":{"ok":true,"rows":[{"name":"vectorscale","default_version":"0.9.0","installed_version":null},{"name…
```

It stops at `{"name…` — the second row of `pg_available_extensions` and a third query's results were both
lost. Not truncated at a record boundary with an ellipsis or a "N more rows" marker, just severed
mid-token.

The obvious fallback — pull the hook's execution logs via `get_hook_logs` — didn't work either, because the
job had deleted the temporary hook as instructed. Cleaning up after itself destroyed the only other copy of
the output.

**Why it matters:** this is the same shape as `BACKLOG.md #4` (no way to retrieve a build job's full
completion log, long reports truncate mid-sentence). It makes jobs unreliable for anything diagnostic. A job
that computes an answer is only useful if the answer survives the trip back, and here roughly half of a
small JSON payload didn't.

In this case it was survivable — the load-bearing row (`pg_extension` returning only `plpgsql`) came through
whole, so the finding stands. But that was luck of ordering, not design.

**Workaround, untested:** have the job write its output into a workspace table instead of the job log, then
read the table. Sidesteps the log path entirely. Otto suggested this and it seems sound, but nothing has
tried it yet.

**Ask:** either raise the output limit meaningfully, or truncate at a record boundary with an explicit
marker so it's obvious data is missing rather than something a reader has to notice from a dangling brace.

Source: Cursor job #98454, DoKnow workspace, 2026-08-05 (`blog/experiments/0033`).
