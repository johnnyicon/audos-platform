---
date: 2026-04-09
priority: 2
status: "not filed"
label: No durable, workspace-pinned instructions for Otto that persist across future sessions and syncs
---

Otto doesn't read any repo file automatically at the start of a session — there's no equivalent of a
`CLAUDE.md`/`AGENTS.md` that Otto itself is guaranteed to load and honor. That means a standing
instruction like "this workspace is single-app, never reintroduce nav chrome" has no durable home; it
has to be re-stated in chat every time, and has no effect on platform-initiated actions (template
refreshes, `[audos-sync]` commits) that don't go through a chat turn at all. Directly related to the
repeated `Desktop.tsx` regressions in `bugs/0031` — even a perfect one-time fix doesn't stick if nothing
enforces it going forward.

Source: `throughline-forge/tmp/audos-feature-request-off-platform-development.md`.
