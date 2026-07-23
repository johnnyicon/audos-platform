---
date: 2026-03-31
area: analytics-api
status: fixed
filed: no
label: "`URLSearchParams is not defined` crash"
---

The analytics-api surface crashed with `URLSearchParams is not defined` in the execution environment it
runs in.

> Fixed by rewriting the affected code with a manual `buildQuery()` helper instead of relying on
> `URLSearchParams`. Resolved 2026-03-31, same day it was raised.
