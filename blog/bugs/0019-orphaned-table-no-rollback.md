---
date: 2026-07-16
area: db-api
status: open
filed: yes
filed_ref: "Priority Support, field-notes workspace, 2026-07-16"
audos_status: pending
label: Failed `workspace_db_create_table` call leaves an orphaned, undroppable table
---

A `CREATE TABLE` call with a foreign key that can't be satisfied fails as expected — but the failure
isn't rolled back. It leaves a real, physical table behind that the platform's own catalog doesn't know
about, permanently blocking any future attempt to use that table name.

> Building a `site_sessions` table with `access_id: uuid` referencing another table's `id` column failed
> with `Failed to create table: foreign key constraint "fk_app_site_sessions_access_id" cannot be
> implemented` — a reasonable rejection, since the platform silently generates every `id` as a `serial`
> integer, never the `uuid` we'd asked for (see #0020). The problem: every subsequent attempt to create
> `site_sessions` — including a corrected version with an `integer` `access_id` — failed with `Failed to
> create table: relation "app_site_sessions" already exists`. The orphaned relation is invisible to
> `db_list_tables`/`db_describe_table` but still physically blocks the name. There is no `drop_table` or
> cleanup tool available. Only workaround found: create the table under a different name entirely
> (`site_sessions_v2`). The original name is permanently unusable until Audos clears it on their end.
