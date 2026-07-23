---
date: 2026-03-31
area: db-api
status: fixed
filed: no
label: "Insert into `dashboard_activity` failed — missing required `title` field"
---

Writes to the `dashboard_activity` table failed silently on missing `title`, with no schema documentation
surfacing the requirement ahead of time.

> Resolved by documenting the schema directly: `title` is a required field on every `dashboard_activity`
> insert.
