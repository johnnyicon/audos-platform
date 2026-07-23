---
date: 2026-07-16
area: db-api
status: inconclusive
label: "Does an explicit schema contract up front beat letting Audos improvise structure?"
---

**Hypothesis:** Does giving Audos's agent-driven build process a complete, explicit data contract
up front — a fully designed `SCHEMA.md` — produce a more reliable result than letting the platform
improvise table structure as it goes?

**Method:** Designed the `content_items`/`site_access`/`site_sessions` schema locally before touching the
platform at all. Following the AGENTS.md finding above (a conventions file isn't read automatically),
pasted the full schema directly into the dispatch prompt itself rather than relying on a file the agent
would need to go find.

**Result: partially confirmed, but confounded by an unrelated platform bug — no clean answer.** The
dispatch honored the pasted schema faithfully for two of the three tables. The third failed, but because
of a genuine, separate platform bug (the platform silently generates every `id` as `serial`, rejecting an
explicitly requested `uuid`) — not because the contract-vs-improvisation question resolved either way.
As far as it was exercised, pasting the schema directly into the prompt worked; the experiment just
never got a clean, unconfounded trial to fully answer the original question.

See `field-notes/ACTIVITY-LOG.md`, schema-design rows; `blog/0013-building-field-notes-in-the-open.md`.
