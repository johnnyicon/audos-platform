---
date: 2026-07-16
area: app-build
status: confirmed
label: Can client-driven sequential hook calls substitute for server-side scheduling?
---

**Hypothesis:** Regardless of whether the scheduler itself is reliable, can an app substitute for it by
having the client call a hook repeatedly and sequentially, paced by the client rather than a server-side
timer?

**Method:** Five sequential client-driven hook calls, timed, checked for rate-limiting, auth friction, or
concurrency issues.

**Result: confirmed pass, clean.** 5/5 succeeded, roughly 1–1.8s each, zero rate-limiting or concurrency
friction. This removes the scheduler as a hard blocker for anything that can tolerate the client staying
open during the work — only genuinely unattended, fire-and-forget background work still needs the
scheduler itself fixed.

See `BACKLOG.md #12`; `CHANGELOG.md`, scheduled-hooks section.
