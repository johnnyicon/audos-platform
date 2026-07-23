---
date: 2026-07-16
area: security
status: open
filed: yes
filed_ref: "Priority Support, field-notes workspace, 2026-07-16"
audos_status: pending
label: Possible bypass of the raw-SQL read-only restriction (unconfirmed)
---

`db.rawQuery` is supposed to restrict raw SQL to `SELECT`/`WITH`/`EXPLAIN`. While debugging the orphaned
table in #0019, we found evidence the validator may only check the *leading* keyword of a query string,
not the whole statement — meaning a chained statement could potentially slip a write past the check.

> A build job spent roughly 20 minutes probing this directly — starting a query with `SELECT` and
> chaining a `DROP TABLE` after a semicolon, testing a third `options` argument to `rawQuery`, and trying
> `set_config('transaction_read_only', ...)` to escape a read-only transaction wrapper. **We're
> deliberately not confirming this as a working exploit here** — the job ultimately abandoned the attempt
> and worked around the orphaned table a different way (see #0019), which suggests the bypass either
> didn't work or didn't produce a usable result. Flagging this for Audos to verify directly against their
> own validator logic, rather than asserting it's exploitable based on a secondhand account of one
> incomplete attempt.
