---
date: 2026-07-17
priority: 1
status: "not filed"
label: "A real config option to make a workspace's own custom auth the sole front door"
---

A workspace owner who builds a fully custom sign-in gate (our `site_access` username/password screen)
currently cannot make it the first thing a visitor sees — Audos's own `EmailGate.tsx` is unconditionally
the signed-out view of every workspace, and a custom gate can only render after it. The only way around
this today is rebuilding `EmailGate.tsx` itself from scratch, which also means giving up the platform's
CRM contact registration and ad-pixel firing on every sign-in — which, for a private or internal tool,
is exactly the point, not a loss.

Ask: a real, documented Developer-panel option (not a code rewrite) to either (a) disable EmailGate
entirely for a workspace and let the app's own signed-out content render first, or (b) let a workspace
register its own component as the signed-out entry view. Priority 1 because the current default — a
private tool's visitor data silently flowing into Audos's own CRM and to Meta/Reddit ad pixels — is a
real, undisclosed behavior most builders would not expect and could not discover without reading the
served source directly, as we had to.
