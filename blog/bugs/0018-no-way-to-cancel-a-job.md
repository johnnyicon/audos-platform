---
date: 2026-07-14
area: app-build
status: open
filed: no
label: No way to check live progress or cancel a running Cursor delegation job
---

No diagnostic or abort tool exists for a running `cursor_delegation` job. `list_jobs` shows status but
its outcome/error field stays empty until the job reaches Complete/Failed; `get_hook_logs` only covers
server-function hooks; `get_audos_code_status` only covers Audos Code threads; the harness's own
`TaskStop`/`TaskList` tools can't see jobs-board tasks at all.

> **Correction:** the specific job that surfaced this (a 5-test capability probe) turned out not to be
> hung — it completed at 45m01s after 67+ real tool actions, genuinely thorough work rather than a stall.
> The underlying gap still stands, though: while a job is running, there is no way to tell a real hang
> apart from slow-but-working, and no way to cancel either one if you need to. A same-scope follow-up job
> we dispatched in the meantime queued silently behind it rather than running in parallel or erroring —
> confirming a related, previously-flagged finding about same-app-scope job concurrency.
