---
date: 2026-07-16
area: app-build
status: inconclusive
label: "Does an AGENTS.md-style conventions file steer Audos's own code generation?"
---

**Hypothesis:** Some agent platforms read a conventions file (AGENTS.md, .cursorrules, CLAUDE.md) and
apply it to their own code generation. Does Audos?

**Method:** Checked the workspace file tree for any such file by default (none found). Behavioral test:
wrote an AGENTS.md with an explicit rule ("all new hook files must have a dated header comment"),
generated a hook via `platform.generateText` with a neutral prompt, checked whether the convention was
honored. A separate direct probe asked the generation endpoint whether it could see any conventions
file at all.

**Result: confirmed fail for the `generateText` path — but explicitly not a full answer.** No
conventions file exists by default; the written convention was ignored; the endpoint responded "NO
WORKSPACE FILE CONTEXT AVAILABLE." Left deliberately open: whether the Cursor-agent/editor-agent flow —
the one that does most of DoKnow's actual app building, and the one that matters most — respects a
conventions file, since testing that requires the file to exist *before* a job starts, not something
discoverable mid-run via `generateText`. The status here reflects that the question that actually
matters for day-to-day building is still unanswered.

See `CHANGELOG.md`, "TEST G" section; `field-notes/ACTIVITY-LOG.md`, row on this same probe.
