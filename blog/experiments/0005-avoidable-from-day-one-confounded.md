---
date: 2026-07-13
area: app-build
status: corrected
label: "\"One Good Thing\": briefing all four shell bugs upfront — confounded by a shared file"
---

**Hypothesis:** If all four known chat-shell bugs (full-screen landing, deep-link resolution,
chat-popup-not-teardown, `JSON.stringify` on writes) are briefed explicitly before any code exists, a
brand-new app avoids all four from birth.

**Method:** Built "One Good Thing" (job #81791) from a brief stating all four fixes upfront, then
verified live by hand — zero-delay screenshot, a real database write, a chat-affordance click.

**Result: reported as a full pass, corrected the next day to "only one-quarter right."** Three of the
four items live in a workspace-shared `Desktop.tsx`, which had already been fixed in this same workspace
days earlier — the new app inherited already-clean shell mechanics rather than testing whether briefing
upfront actually prevents the bugs. Only the `JSON.stringify`-on-write item was a genuine, unconfounded,
per-app test, and that one held. The correction mattered enough to reshape how later tests were designed
— see the "Red Pill" experiment, which deliberately tested a zero-briefing baseline to isolate what's
workspace-inherited from what's actually per-app.

See `docs/platform/22-eliminating-the-chat-shell-playbook.md`, "'Avoidable from day one' — original
claim and a correction," and `blog/0010-building-it-right-the-first-time.md`.
