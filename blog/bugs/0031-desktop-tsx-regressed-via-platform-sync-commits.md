---
date: 2026-04-09
area: app-build
status: open
filed: no
label: Desktop.tsx got silently reverted by a platform-initiated commit, at least the third time the same fix has been undone
---

A deliberate customization — stripping the multi-app sidebar and mobile bottom-nav out of `Desktop.tsx`
so the app renders full-canvas, single-app — had held for months after being properly deployed. On
April 8–9, an `[audos-sync]` commit (`c326891`) rewrote `Desktop.tsx` and silently reintroduced the
multi-app sidebar, with hardcoded `coreApps` entries for apps (Signature/Studio/Briefing) that weren't
even registered in the workspace's own `config.json` at the time. Production UI broke as a result.

> **This was at least the third time the same fix had been undone**, and root-causing it hit a wall
> that's itself worth recording: there's no audit trail on platform-initiated `[audos-sync]` commits —
> no session ID, no prompt pointer, no "who/why." Asked Otto to investigate; its best guess was that it
> had made the edit during a previous session, but that could not be confirmed from either side. The
> commit diff is visible. Nothing about its origin is.

Two structural gaps make this close to inevitable for any team working IDE-first against an Audos
workspace: no file-preservation mechanism (no way to mark `Desktop.tsx` "don't regenerate"), and no
config-driven layout (no field in `config.json` to declare the space's shape, so the only way to run
full-canvas is imperative JSX surgery that any future refresh can silently undo). Both are filed as
their own feature requests. In the meantime, stopgaps only — a GitHub Action tripwire that fails CI on
known regression patterns, and reverting back to clean state by hand each time it recurs. Detection,
not prevention.

Source: `throughline-forge/tmp/audos-feature-request-off-platform-development.md`.
