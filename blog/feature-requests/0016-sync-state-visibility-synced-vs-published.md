---
date: 2026-04-17
priority: 1
status: "not filed"
label: "\"Synced\" and \"published\" are different states, and the UI only ever shows one of them"
---

The Sync Activity panel marks a GitHub push "done" the moment source lands in the workspace — but the
actual compile step (source → `.published-source/`) runs separately, asynchronously, on an undocumented
~15–30 minute latency window. From the outside, "Sync Activity says done" and "the site is still serving
the old bundle" look identical to "the pipeline is broken," and it takes real investigation to tell them
apart.

> **Worth telling straight, because the story itself makes the case.** A first read concluded the
> GitHub-sync path "never compiles app code" — the published bundle still matched the old version hours
> after a push. A second, more patient look the next morning found the opposite: source and published
> hashes matched, the new code was live, the first check had simply landed inside the compile window
> before it finished. The draft got retracted before it was ever sent, with a note admitting exactly
> that: "I'd have sent it without the second look if I hadn't double-checked." The pipeline wasn't
> broken. The *visibility* into where a push currently stood was.

Two small asks would have prevented the whole misdiagnosis: a second state in the Sync Activity row
(`Synced • Compiling…` → `Synced • Published`), or simply surfacing the `.published-source/` hash next
to the source hash so a match/mismatch is visible at a glance. Two adjacent asks from the same
investigation: document the expected sync→compile→CDN latency so "normal" has a number attached to it,
and document what switching between GitHub Dev Mode and Platform Dev Mode actually does to existing
source, drafts, and the GitHub connection — asked directly, Otto had no documented answer for any of it.

Source: `throughline-forge/tmp/audos-feature-request-followup.md`.
