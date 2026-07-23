---
date: 2026-04-20
area: identity
status: confirmed
label: useWorkspaceDB and useSession are Audos-hosted-app only — they don't work for a CDN-ported external app
---

**Hypothesis:** Since `useWorkspaceDB` and `useSession` are documented as the platform's React data/auth
hooks, they should work the same way regardless of where the compiled app is actually served from —
Audos's own CDN, or an external host after porting the code out.

**Method:** Attempted to keep using these hooks after porting Throughline's apps off Audos-hosted
deployment onto an external CDN/Railway-style setup, to avoid rewriting the data layer twice.

**Result:** They don't travel. Both hooks are Audos-hosted-app only — confirmed directly in the team's
own service documentation ("Not Used... Audos-hosted-app only, not applicable on Railway"). A CDN-served
external app has to go back to plain REST calls against `db-api` directly, the same pattern this
project's own `blog/0002` already committed to for a different reason (avoiding `useWorkspaceDB`'s
silent `shared: false` session-scoping default). Separately unresolved: whether `useSession()` is even
importable at all outside an Audos-hosted app, or errors immediately on import — noted as an open
question, not confirmed either way.

Practical takeaway: any app planning to eventually run outside Audos-hosted deployment shouldn't build
against these two hooks in the first place — the REST/`db-api` path this project already prefers is also
the only one that's portable.

Source: `throughline/docs/audos-services.md`, `throughline-forge/tmp/otto-identity-auth-prompt.md`.
