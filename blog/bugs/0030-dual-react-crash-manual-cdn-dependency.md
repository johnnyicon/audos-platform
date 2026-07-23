---
date: 2026-04-16
area: app-build
status: open
filed: no
label: Adding a React-dependent CDN package Audos doesn't pre-configure silently loads a second React and crashes (React error #31)
---

Audos auto-generates a `?deps=react@18.3.1`-style pin in the esm.sh importmap URL, but only for its own
pre-configured packages (`react`, `react-dom`, `lucide-react`, `react-markdown`). Adding any other
React-dependent package by hand — in this case `@tanstack/react-router` — gets a bare esm.sh URL with no
such pin. esm.sh then resolves that package's own React peer-dependency independently, loading React 19
alongside the platform's React 18. Result: a blank white screen and React error #31 (two React copies
in the same tree), with no obvious cause from the importmap itself — the URL looks correct on inspection.
Only visible by checking network requests for two separate React bundles loading.

Confirmed by directly diagnosing and removing the offending import — but that alone didn't immediately
fix it, because Audos's ESBuild pipeline compiles every `.tsx` file in the app directory, not just the
entry point's import tree (see the matching experiment). A sibling file still importing the same
library kept the crash alive until every file in the directory was cleaned, not just `App.tsx`.

Worked around by replacing the router with a ~30-line custom state router with no CDN dependency at
all, avoiding the problem rather than fixing it. The underlying gap — no way to tell Audos to add a
`deps` pin for a manually-added package — remains; see the matching feature request for the concrete
ask (a `deps` field in `cdnDependencies`, or documentation that any React-dependent package added by
hand needs the same treatment).

Source: `throughline/docs/working/audos-migration-plan.md`, `throughline-forge/tmp/nick-throughline-on-audos.md`, `nick-routing-platform-requests.md`.
