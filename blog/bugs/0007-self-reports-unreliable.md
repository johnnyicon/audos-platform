---
date: 2026-07-12
area: app-build
status: open
filed: no
label: Job self-reports don't reliably match live reality
---

Jobs have reported "complete," "published," and even explicit "verified" success language that did not
match the live app when independently re-tested in-browser — not once, but repeatedly, across unrelated
apps and fixes. No malicious pattern; per Otto's own analysis this is likely an **authority-boundary
gap**: the Cursor agent reports sincerely on what it did inside its own sandbox, while the steps that
determine actual ground truth (confirmation gates, publish/recompile, what's actually served) are
enforced Audos-side, outside the agent's own visibility.

> Standing practice adopted because of this: never trust a job's self-reported status — always
> independently verify the live, served result. See the related feature requests for an automated
> post-build smoke check and a real publish-status endpoint, either of which would close this gap at the
> source instead of requiring a human re-check every time.
