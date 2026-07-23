---
date: 2026-04-09
priority: 1
status: "not filed"
label: A way to mark a file "user-owned, don't regenerate" — no pragma, no .audosignore, nothing exists today
---

There's no mechanism today to tell Audos "this file is mine, never overwrite it" — no pragma comment,
no `.audosignore`-equivalent manifest flag, nothing. Whenever Otto, a platform template refresh, or a
`[audos-sync]` commit touches a file, hand-made customizations get silently overwritten with no warning
and no diff review. This is the direct, confirmed cause of a real production regression (see
`bugs/0031`) — a deliberately customized `Desktop.tsx` got reverted at least three times by exactly this
mechanism.

The concrete ask, as simple as it sounds, would resolve most of what triggers this class of finding:

```
# .audosignore
Desktop.tsx
config.json
```

Otto itself, asked directly about the gaps that make these regressions close to inevitable for
IDE-first teams, named this as the single highest-priority fix if it were making the call.

Source: `throughline-forge/tmp/audos-feature-request-off-platform-development.md`.
