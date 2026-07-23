---
date: 2026-07-16
area: app-build
status: confirmed
label: "field-notes' core experiment: can a genuinely fresh workspace be born outside the chat shell?"
---

**Hypothesis:** Every prior "shell-escape" result (the library probe, the base-URL landing config, the
Red Pill migration) ran in a workspace whose shared `Desktop.tsx` had already been fixed once. Can a
genuinely fresh Audos workspace — never touched by any prior fix — be born full-screen, outside the chat
shell, from the very first paint, purely by briefing the 4-item shell-escape checklist upfront, with no
retrofit required? This is the harder version of the earlier experiments, run in a workspace with
nothing to inherit.

**Method:** Created a genuinely new workspace (`field-notes`). First attempt briefed the full checklist
and 5-point verification upfront (job #83055). It crashed mid-run on unplanned cosmetic polish before
posting a completion report; checking the raw file-tree API afterward came back empty, initially read as
data loss. A scoped retry (job #83114) was told the database work was already done and instructed to
stop immediately once the 5 checks passed. It completed cleanly, reporting the previous job's work was
already present. Verification was then done independently by hand: the workspace's own Preview panel,
Draft compared against Live side by side, logging in with test credentials personally, reaching the real
content feed, clicking the chat affordance personally.

**Result: initially "inconclusive, not proven either way," resolved the same day to confirmed YES,
independently verified.** The "inconclusive" verdict turned out to rest on a mistake in the verification
method itself — the file-tree endpoint checked only reflects the *published* bundle, never a draft. Once
checked correctly (Draft vs. Live in the Preview panel), the app had been sitting there the whole time:
full-screen from first paint, no chat-shell flash, no retrofit needed. This also triggered a retraction
of a "file writes not durable" bug report as a false alarm — though a later, more careful look found even
that retraction had rested on an unverified assumption about the intermediate state between the two jobs,
which is now permanently unrecoverable to check further.

See `docs/platform/22-eliminating-the-chat-shell-playbook.md`; `blog/0013-building-field-notes-in-the-
open.md`; `field-notes/ACTIVITY-LOG.md`.
