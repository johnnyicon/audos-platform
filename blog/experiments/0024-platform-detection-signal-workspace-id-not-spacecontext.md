---
date: 2026-04-16
area: identity
status: corrected
label: The right way to detect "am I running inside Audos at all" is window.__WORKSPACE_ID__, not window.__spaceContext
---

**Hypothesis:** `window.__spaceContext` — the object `blog/0002`'s multi-tenant post already established
as the accessor for the *current signed-in user's* identity (`__spaceContext?.username`, populated once
EmailGate completes) — would also be the right signal for a different, narrower question: not "who is
signed in," but "is this code running inside an Audos workspace at all," for an app built to run both
on-platform and ported elsewhere (Railway, a CDN).

**Method:** Built sign-in gating logic against `window.__spaceContext` truthiness as the on-platform
check, then watched real behavior in production.

**Result:** Wrong, in a way that mattered. `window.__spaceContext` as a bare platform-presence flag is
never actually set in the context this was tested from — the app showed its own sign-in gate to users
who were already authenticated on the platform, because the detection check itself never passed.

> **Correction, found empirically:** the reliable, always-injected global for "is this Audos" is
> `window.__WORKSPACE_ID__`. Switched the check to `const isOnPlatform = () => !!(window as any).__WORKSPACE_ID__`
> and it worked immediately.

Worth being precise about what this does and doesn't correct: `blog/0002`'s claim about
`__spaceContext?.username` as the *user-identity* accessor stands — that's a different property, used
for a different purpose (who's signed in, not whether the platform is present at all), and nothing here
contradicts it. What's new is a second, distinct signal this project hadn't needed before: reliable
platform-presence detection for code meant to run both on and off Audos.

Source: `throughline-forge/tmp/nick-throughline-on-audos.md`.
