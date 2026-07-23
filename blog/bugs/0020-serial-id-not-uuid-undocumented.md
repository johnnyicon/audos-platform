---
date: 2026-07-16
area: db-api
status: open
filed: yes
filed_ref: "Priority Support, field-notes workspace, 2026-07-16"
audos_status: pending
label: Table `id` is always a platform-generated `serial`, never the `uuid` you ask for — undocumented
---

`workspace_db_create_table` reserves the `id` column for itself. Asking for `id: uuid` doesn't get a
uuid primary key — it gets rejected outright, and the platform's own reference docs never say so.

> Requesting `id: uuid` in a column list returns `column "id" specified more than once` — the platform
> always injects its own auto-increment `serial` integer `id`, and won't let a caller override its type.
> `uuid` remains valid as a non-primary-key column. Also undocumented: every table silently gets two
> extra columns beyond what's requested — `session_id` (text) and `updated_at` (timestamp, default
> `NOW()`). None of this appears in the current column-type reference, which lists `uuid` as simply one
> of the supported types with no caveat about primary keys. Practical consequence: any schema designed
> assuming a `uuid` primary key needs to design foreign keys as `integer` instead, and needs to expect
> two uninvited columns on every table.
