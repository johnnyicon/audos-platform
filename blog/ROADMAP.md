# Blog/site roadmap — deferred, spec'd now so nothing gets lost

Captured 2026-07-13. **Not being built yet** — explicitly deferred by request. This is the spec to
execute against when picked back up, ideally via a sub-agent doing the review/authoring pass described
below.

## What exists today (confirmed working)

The scorecard grid at the top of the blog page is still there — 9 cells, one per entry, each showing a
status chip (`pass` / `fixed` / `open`). What it does **not** do: give a distinct count of bugs found
per entry, or link out to a dedicated bug record. That's the actual gap driving this roadmap.

## The target structure: three pages, one shared nav

1. **Blog** (exists) — the narrative entries, as now.
2. **Bugs** (new) — one row/card per bug, pulled from `BACKLOG.md` and
   `docs/platform/bug-reports/*.md`, each showing: description, **fixed or still open**, **whether it's
   been filed with Audos's own help desk** (and if so, the report file/link), and current status. This
   is a re-render of data that already exists in `BACKLOG.md` — no new data model needed, just a second
   page generated from the same source.
3. **Feature Requests** (new) — currently all bundled together as a list inside
   `docs/platform/20-power-user-wishlist.md`. Needs to become **one file per feature request** (a new
   `feature-requests/` directory, same frontmatter convention as blog posts: title, rationale, status),
   so each can be linked to individually and the page can render them as distinct cards.

## The authoring pass needed before this can be built

A sub-agent (when this is picked back up) should:
1. Read `doknow-kb/audos/ACTIVITY-LOG.md` (a different repo — see `docs/platform/24-where-new-findings-go.md`), `BACKLOG.md`, all of `docs/platform/*`, and all blog posts.
2. Extract every distinct bug already found (fixed or open) and confirm/complete its status fields
   (fixed/open, filed-with-Audos yes/no + link, current status) — most of this data already exists in
   `BACKLOG.md`, this is a restructuring pass, not new research.
3. Extract every distinct feature-request-shaped idea (most already live in
   `docs/platform/20-power-user-wishlist.md`) and split each into its own file under `feature-requests/`.
4. Extend `scripts/build_blog.py` (or add sibling scripts) to generate `bugs.html` and
   `feature-requests.html` alongside `dist/index.html`, sharing the same design tokens/masthead, with a
   simple top-level nav linking all three.

## Why deferred rather than done now

Explicitly requested to hold this for later. Capturing it here in full, rather than as a half-remembered
verbal note, so a future session or sub-agent can execute directly from this file without re-deriving
what was meant.
