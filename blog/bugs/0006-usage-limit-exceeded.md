---
date: 2026-07-12
area: app-build
status: open
filed: no
label: Build jobs instant-fail with usage_limit_exceeded (Audos's own account, not the workspace owner's)
---

Cursor Background Agent jobs can instant-fail at the delegation step with `usage_limit_exceeded`. The
error links to `cursor.com`, but the account is almost certainly **Audos's shared build infrastructure**,
not the workspace owner's — so the owner usually can't clear it self-serve. The failure is account-level
and typically transient: jobs that instant-failed in a session succeeded on retry a short time later.

> No user-facing status or ETA for this exists today. Recommendation if it persists across retries over
> time: file with Audos as a platform-capacity issue, with exact error text, task IDs, and timestamps —
> do not tell a user to go raise a Cursor account limit they don't own.
