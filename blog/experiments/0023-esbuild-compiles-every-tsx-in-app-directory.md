---
date: 2026-04-16
area: app-build
status: confirmed
label: Audos's ESBuild pipeline compiles every .tsx file in the app directory, not just the entry point's import tree
---

**Hypothesis:** Audos's build pipeline compiles `App.tsx` and whatever it actually imports, the way any
normal bundler treats an entry point — a sibling file nobody imports shouldn't affect the runtime bundle.

**Method:** While diagnosing a dual-React crash (see the matching bug report), stripped a bad import
(`@tanstack/react-router`) out of `App.tsx` and redeployed, expecting the crash to clear immediately.

**Result:** It didn't. Another file in the same `apps/{id}/` directory — one `App.tsx` never imports —
still had the same import, and that was enough for the library to load at runtime and reproduce the
crash. Only after cleaning up every sibling file's imports, not just the entry point's, did it resolve.

> The actual rule: Audos's ESBuild compiles **every `.tsx` file in the app directory**, regardless of
> whether the entry point's import graph ever reaches it. This isn't necessarily wrong behavior — just
> undocumented and easy to get bitten by, since it breaks the normal mental model of "unused imports in
> unreferenced files are dead code." Once you know the rule — everything in `apps/{id}/` is live, full
> stop — you plan around it by keeping unused/experimental files out of that directory entirely rather
> than assuming they're inert.

Source: `throughline-forge/tmp/nick-throughline-on-audos.md`.
