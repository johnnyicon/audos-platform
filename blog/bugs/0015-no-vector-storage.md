---
date: 2026-07-14
area: db-api
status: open
filed: no
label: No native vector/embedding storage, and no way to enable pgvector
---

`workspace_db_create_table` rejects a `vector(N)` column type — table creation is restricted to a fixed
enum (`text/integer/bigint/decimal/boolean/timestamp/date/json/uuid`) with no vector type in it. Raw SQL
is `SELECT`/`WITH`/`EXPLAIN`-only, so `CREATE EXTENSION vector` can't be run either.

> `pg_available_extensions` confirms pgvector 0.8.1 is installed on the Postgres binary
> (`installed_version: null`) — it's sitting right there, just not enabled, with no path found to enable
> it ourselves. Only workaround: store embeddings as a JSON array and brute-force similarity search in
> hook JavaScript — workable for a handful of items, not a real retrieval index. This is the single
> biggest blocker for any retrieval-grounded generation pipeline on this platform.

**Re-verified 2026-08-05, three weeks on — still holds.** `pg_extension` returns exactly one row,
`plpgsql`; `vector` is not installed. `vectorscale` 0.9.0 still shows `installed_version: null`.
Getting that answer took four routes (Otto can't run raw SQL, no hooks exist on the workspace, the
external Otto can't create one) and finally a Cursor job — see `blog/experiments/0033`.

One proportion note the original entry overstated: calling this "the single biggest blocker for any
retrieval-grounded generation pipeline" was too strong. `BACKLOG.md #13` later benchmarked the
JSON-array + brute-force fallback at real embedding size (1,536 floats, 1–3ms at 1,000 rows). It's a
ceiling on headroom at larger scale, not something blocking a build today.
