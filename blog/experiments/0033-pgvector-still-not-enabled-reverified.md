---
date: 2026-08-05
area: db-api
status: confirmed
label: Re-verified three weeks on — pgvector is still not enabled, and getting a straight answer took three failed routes first
---

**Hypothesis:** `bugs/0015` established on 2026-07-16 that pgvector sits in Audos's extension catalogue
but is never actually installed. Three weeks later, before repeating that claim to Audos directly, it
needed re-checking — the platform moves fast, and a stale finding stated confidently is worse than no
finding at all.

**Method:** what's interesting here is how hard a two-line SQL query turned out to be. Four routes,
three dead ends:

1. **Ask Otto to run the SQL.** Refused, correctly — its `db_query` tool is a *structured* interface
   (table name, filters, sort), takes no raw SQL, and cannot reach system catalogue tables at all. To
   its credit it said so plainly rather than fabricating rows.
2. **Call an existing server function** that could run `db.rawQuery` (which does permit `SELECT`).
   Probed five plausible hook names on the workspace — `db-api`, `database`, `db`, `query`,
   `raw-query` — all returned `404 Hook not found`. None exist.
3. **Have Otto create a diagnostic hook.** Otto reported it has no `manage_server_functions` tool at
   all — which **directly contradicts `docs/platform/07`**, where that tool is documented as *the* way
   server functions get created. See the separate finding below.
4. **Stage a Cursor delegation job** — the route that worked. A job (model pinned to `fable-5`) created
   a temporary read-only `ext-check-diag` hook, ran three catalogue `SELECT`s, logged the raw JSON, and
   deleted itself. Dispatched by naming the specific draft ID rather than a bare "run it," per the
   standing rule that a generic confirmation can execute every pending draft.

**Result: confirmed, the July finding holds.** `pg_extension` returned exactly one row — `plpgsql` —
complete and untruncated. **`vector` is not installed.** `vectorscale` 0.9.0 appears in
`pg_available_extensions` with `installed_version: null`, consistent with July.

One honest gap: the jobs-list tool truncated the output mid-JSON, cutting the `vector` row of
`pg_available_extensions`. So "not enabled" is verified as of today; "present in the catalogue" rests on
July's evidence plus the `vectorscale` row pointing the same way. The installed list is the load-bearing
half and it came through whole.

**Worth keeping in proportion.** `BACKLOG.md #13` already established that the JSON-array plus
brute-force cosine fallback benchmarks fine at realistic dimensionality (1,536 floats, 1–3ms at 1,000
rows). So this is a ceiling on headroom at larger scale, not something blocking a build today — and it's
worth saying that out loud rather than letting a real gap get inflated into a blocker it isn't.

**The meta-finding:** answering "is this extension installed" — about thirty characters of SQL — required
dispatching a background build job to a third-party agent runner. That is the friction this whole log
keeps documenting, in miniature.

Source: Otto chat + Cursor job #98454, DoKnow workspace (`8a65a4ac-5a22-435f-b55f-c41ea34ca00d`),
2026-08-05. Logged in `doknow-kb`'s `audos/ACTIVITY-LOG.md`.
