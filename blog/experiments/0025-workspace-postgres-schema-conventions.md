---
date: 2026-04-17
area: db-api
status: confirmed
label: Reverse-engineering the workspace Postgres schema conventions end to end
---

**Hypothesis:** Audos workspace databases follow a consistent, discoverable naming and ownership
convention — worth mapping once, properly, rather than re-deriving it piecemeal on every project.

**Method:** Directly inspected a live workspace's Postgres schema (table names, columns, generic hook
behavior) rather than relying on SDK docs alone, and classified every table by who actually owns it.

**Result:** Confirmed a consistent pattern:

> Schema name follows `ws_<workspace-uuid-with-underscores>`. Every platform/app table carries an
> `app_` prefix. Most tables also carry `session_id`/`user_id`/`org_id` columns. Audos's own platform
> infrastructure (`funnel_contacts`, `ad_campaigns`, `workspaces`, etc.) lives in a **separate platform
> database**, entirely outside the workspace schema — so workspace-schema changes can't touch platform
> internals. Per-app scaffold tables (this project's Briefing/Guest-Prep/Studio/Podcast-Setup/Reels
> tables) are safely droppable once an app's own backend takes over data ownership. The generic REST
> hooks (`db-api`, `ai-api`, `email-api`) work against *any* workspace table without additional
> per-table endpoint setup — `describe`/`query`/`insert`/`update`/`delete` actions, and update/delete
> both require at least one filter condition as a basic safety rail.

This consolidates and confirms several scattered observations from this project's own earlier
DoKnow/field-notes work (the `ws_<uuid>` schema convention, the `app_` prefix, the platform-generated
`id`/`session_id`/`updated_at` columns) into one clean reference, rather than something newly
discovered — but it's the first time it's been written down as a single, complete picture end to end.

Source: `throughline/docs/audos-throughline-api-doc.md`, `throughline-forge/tmp/otto-audos-db-conventions-prompt.md`.
