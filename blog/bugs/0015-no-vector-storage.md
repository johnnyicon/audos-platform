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
