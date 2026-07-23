---
date: 2026-04-18
area: app-build
status: open
filed: no
label: Desktop.tsx's hash router matches the entire hash, so any sub-path breaks — no deep links for any third-party app
---

The platform shell's hash-based router (`Desktop.tsx`) identifies which app to render by matching the
**entire hash string** against a registered app ID — not a prefix. A URL like `#throughline/episodes`
doesn't match the registered app `throughline`, so it falls through to whatever the default app is
instead of rendering the intended route. Practical effect: no deep links, no working back/forward
history, no bookmarkable URLs for any sub-path of a third-party app — every app is confined to exactly
one hash value, full stop.

A one-line fix was identified and supplied informally (matching on `hash.split('/')[0]` instead of the
full string), but whether Audos adopted it wasn't confirmed — this is the workaround supplied, not a
verified platform fix. Filed as a routing capability request rather than a formal bug report.

Source: `throughline-forge/tmp/nick-routing-platform-requests.md`.
