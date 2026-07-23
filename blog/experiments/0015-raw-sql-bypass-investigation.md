---
date: 2026-07-16
area: security
status: inconclusive
label: "Does db.rawQuery's SELECT/WITH/EXPLAIN restriction check only the leading keyword?"
---

**Hypothesis:** `db.rawQuery` is supposed to restrict raw SQL to `SELECT`/`WITH`/`EXPLAIN`. Does its
validator check only the leading keyword of a query string, meaning a chained statement (`SELECT 1;
DROP TABLE ...`) could slip a write or DDL statement past the check?

**Method:** Roughly 20 minutes of live testing from multiple angles while trying to clear an orphaned
table: transaction-scoped read-only wrappers, `set_config('transaction_read_only', ...)`, a third
`options` argument to `rawQuery`.

**Result: inconclusive, explicitly unconfirmed either way.** The job abandoned this approach and worked
around the orphaned table a different way instead (creating a differently-named table) rather than
continuing to push on the bypass. That's circumstantial evidence the bypass didn't actually work, but
it's not proof — the attempt was incomplete, not a clean negative result. Filed for Audos to verify
directly against their own validator logic, not asserted as a working exploit.

See `BACKLOG.md #17`.
