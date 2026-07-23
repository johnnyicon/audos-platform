---
date: 2026-07-13
area: app-build
status: confirmed
label: "\"Red Pill\": migrating an existing app in place, and isolating what's actually per-app"
---

**Hypothesis:** Can an already-platform-generated app be pulled out of the chat shell in place — same
app id, same file, same data — rather than rebuilt as a new app and cut over? And, following the
correction above: are shell mechanics actually workspace-shared, or per-app?

**Method:** Recreated a plain describe-then-generate baseline app (`doknow-home`, job #81930) with
deliberately zero special briefing, expecting it to reproduce the shell-flash bug. Then migrated it in
place with a verbatim design push (job #81943, corrected after a first attempt only pasted CSS literally
without wiring it up) while preserving existing database wiring. Checked four things independently:
mechanics, design fidelity, data continuity, and scope (an untouched sibling app).

**Result: confirmed pass on all four checks — and it disproved the premise it was testing.** The
zero-briefing baseline app rendered full-screen with no flash anyway, proving shell mechanics are
inherited workspace-wide via the shared `Desktop.tsx`, not something each new app needs to be taught
individually. The verbatim design port plus data-continuity migration succeeded, independently verified
live. This result is what triggered the retroactive correction of the "One Good Thing" experiment above.

See `docs/platform/22-eliminating-the-chat-shell-playbook.md`, "Migrating an EXISTING app in place, not
replacing it," and `blog/0011-taking-the-red-pill.md`.
