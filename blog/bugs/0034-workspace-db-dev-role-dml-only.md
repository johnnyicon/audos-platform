---
date: 2026-04-19
area: db-api
status: open
filed: no
label: The workspace's own direct-Postgres dev role is DML-only — every schema change still has to go through Otto in chat
---

Direct Postgres credentials for a workspace (`ws_dev_*` role) grant read/write on existing rows, but all
DDL is blocked: `CREATE TABLE` returns `permission denied for schema`, and both `DROP TABLE` and
`ALTER TABLE` return `must be owner of table`. In practice this means the credentials feature — sold as
direct database access — can't actually do schema management. Any new table, dropped table, or column
change still has to be requested through Otto in chat, the exact bottleneck direct DB access was meant
to remove.

For a team running a daemon that needs to manage its own schema programmatically, this is a real gap,
not a cosmetic one — chat-mediated schema changes aren't viable for automated, ongoing schema
management. The ask was straightforward: admin-level access, so schema changes don't require going
through Otto at all. No confirmation captured on whether that was granted.

Source: `throughline-forge/tmp/nicholas-db-credentials-update.md`.
