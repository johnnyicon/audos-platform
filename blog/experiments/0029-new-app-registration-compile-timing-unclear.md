---
date: 2026-04-17
area: app-build
status: inconclusive
label: Does newly-registering an app via GitHub push skip the compile queue, or just take longer than expected?
---

**Hypothesis:** Adding a second app entry to `config.json` (pointing at a new component file) via a
GitHub push should compile the same way an update to an existing app does — land in `.published-source/`
within the normal sync latency window.

**Method:** Pushed a new app entry (`throughline-next`) via GitHub Dev Mode, confirmed the DB app record
got created, then checked `.published-source/apps/throughline/` for the new component.

**Result:** The new component never appeared in `.published-source/` — only the original app's existing
files were there. Otto's own read: "GitHub sync writes files and updates `config.json`, but does NOT
compile new app components into `.published-source/`... new component registration requires a separate
compile trigger" (`delegate_app_edit` or the Developer panel UI).

> **Marked inconclusive on purpose, not confirmed.** The same investigation had already produced one
> false alarm minutes earlier: an *update* to an existing app looked uncompiled too, until a more patient
> re-check found the compile just hadn't finished yet (~15–30 min sync→compile→CDN latency, undocumented
> but real). Whether the new-app case is a genuinely different code path that skips the compile queue, or
> simply the same latency window not waited out long enough, was never actually settled — the original
> report says so explicitly: "Unclear whether this is the same latency pattern... or whether new-app
> registration really does skip the compile queue." Worth someone re-testing with a longer wait before
> treating this as a confirmed platform gap.

Source: `throughline-forge/tmp/audos-feature-request-followup.md`.
