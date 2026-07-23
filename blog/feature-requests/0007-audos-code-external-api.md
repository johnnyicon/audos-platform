---
date: 2026-07-13
priority: 3
status: not filed
label: Let API-driven agents use the Audos Code backend, not just Cursor
---

When Cursor is over Audos's own shared usage limit, an API-only agent has no working build path today —
it can create drafts, but a human must run them from inside a signed-in browser session. Audos Code exists
specifically to sidestep the Cursor limit, but it requires a signed-in workspace session's user email for
attribution and cannot be launched from the external onboarding/chat API — both the draft-run and
direct-handoff routes fail with "no user email for attribution." This was the one case in this whole build
where we were fully blocked with no workaround at all except waiting.
