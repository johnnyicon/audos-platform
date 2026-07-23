---
date: 2026-04-20
product: Throughline
status: pass
label: We ported Throughline onto Audos, made it work, then mostly left anyway
---

# We ported Throughline onto Audos, made it work, then mostly left anyway

Everything in this log so far about Throughline (`0001`, `0002`) comes from the earliest days — a
platform shell bug, a multi-tenant auth pattern — both true, both narrow. There's a bigger arc sitting
underneath them that never made it into this record: a real, disciplined attempt to port Throughline's
actual product onto Audos-hosted deployment, followed by mostly leaving anyway. Both halves are worth
telling, because neither one alone is the honest version.

## The port was real, not a toy

This wasn't a proof of concept. The ground rules going in were strict: copy the codebase rather than
move it, so the original stayed a working hedge; write tests before the ported code, not after; verify
feature parity empirically before flipping any switch; no propagated debt carried forward just because
it was easier to copy than to fix. Twelve lib modules, a full route tree, an app shell, twelve page
stubs — 23 test files and 321 tests by the time the core shell was live, zero regressions in the
original codebase the whole way through.

It hit exactly the platform edges this log has spent the rest of its time cataloguing. A router spike
using TanStack Router via esm.sh looked fine in isolation and then caused a dual-React crash once
deployed for real — Audos auto-pins React only for its own pre-configured packages, and a manually-added
router pulled in its own second copy silently. The fix was to throw the dependency away entirely: a
30-line hand-rolled state router, no CDN import, problem gone. A platform-detection check built against
`window.__spaceContext` turned out to be checking a global that's never actually set for that purpose —
`window.__WORKSPACE_ID__` was the real, always-injected signal, found only after the wrong check had
already been shipped and quietly shown the app's own sign-in gate to users who were already
authenticated on the platform. Small, specific, exactly the shape of finding this whole SDK has been
built to catch.

> By Phase 4, it wasn't a demo anymore. `App.tsx` was the real multi-file entry point, `config.json`
> pointed at it, and it was live, serving actual traffic, on `app.trythroughline.com`.

## And then the real deployment went to Railway anyway

Today, the actual production app is a Vite build shipped through Railway, backed by Neon Postgres and
Brevo for email. The `audos-workspace/` tree — the one all that porting work went into — is now
explicitly labeled legacy platform material in the team's own deployment docs, not to be touched for
normal product changes.

What stayed, deliberately, is narrower and more specific than "we left":

> Audos Auth → replaced by self-owned OTP/JWT. Audos Database → replaced by Neon. Audos Email →
> replaced by Brevo. Audos's own React hooks (`useWorkspaceDB`, `useSession`) → confirmed Audos-hosted-app
> only, don't work once code is running somewhere else. What's actually still in active use: AI
> generation, called over plain HTTP from the daemon, with the API URL as an env var and an explicit
> "swap the endpoint to any OpenAI-compatible provider without touching code" exit ramp built in from day
> one. Everything else on the "keep" list — Web scraping, Storage, Scheduler, Analytics, CRM — is marked
> as *not yet wired up*, prioritized for later, not committed to.

That's the actual shape of the decision, and it's a more useful one than either extreme story would be.
Not "Audos couldn't do it" — the port worked, the app really did run there. Not "we're all-in on Audos"
either. What's left is Audos used for exactly the things it's good at from outside its own hosted-app
model — a services layer reached over HTTP, replaceable in one line if the terms ever stop making sense
— while auth, data, deployment, and the product's own React runtime live somewhere the team fully
controls. If field-notes' own struggles this month with EmailGate and publish reliability are any guide,
that's not a coincidence. It's the same lesson arriving from a different direction, a few months
earlier.

Source: `throughline/docs/working/audos-migration-plan.md`, `throughline/docs/audos-services.md`,
`throughline/docs/deployment.md`.
