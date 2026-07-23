---
date: 2026-07-12
area: app-build
status: open
filed: no
label: No way to retrieve a build job's full, untruncated completion log
---

Job completion reports repeatedly cut off mid-sentence, at least eight separate times across this build,
almost always at the exact critical detail — the root-cause line, the fix diff, the model-used
confirmation, the publish line. No tool or API call retrieves the untruncated tail.

> Root-caused (not an Audos withholding choice): the underlying Cursor Cloud Agents API truncates its own
> output stream by design, with no documented full-log endpoint — see `docs/platform/21-cursor-backend-research.md`.
> Doesn't remove the ask: Audos could still persist and expose the full run transcript server-side,
> independent of what Cursor's own stream truncates. Workaround adopted in the meantime: verify live in
> the browser rather than chase the report text.
