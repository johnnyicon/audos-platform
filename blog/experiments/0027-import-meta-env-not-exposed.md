---
date: 2026-04-20
area: app-build
status: confirmed
label: Audos's ESBuild pipeline doesn't expose import.meta.env — Vite-only dev fixtures can't be ported as-is
---

**Hypothesis:** Since Audos's build pipeline compiles standard React/TypeScript, common Vite conventions
like `import.meta.env` for build-time environment flags would carry over unchanged when porting an app
built against Vite locally onto an Audos workspace.

**Method:** Attempted to port a Vite-built app (using `import.meta.env.DEV`-style flags for local dev
fixtures) onto an Audos workspace build target without modification.

**Result:** `import.meta.env` isn't exposed by Audos's ESBuild pipeline — it's a Vite-specific global,
and Audos's bundler doesn't provide an equivalent. Any dev-only fixture/flag gating built against it
silently breaks (or fails to compile) once ported, and has to be rewritten against whatever mechanism
Audos's runtime actually exposes (a runtime constant, a platform-detection check like `window.__WORKSPACE_ID__`
from the platform-detection experiment above) rather than a build-time env substitution.

Source: `throughline/docs/working/audos-migration-plan.md`.
