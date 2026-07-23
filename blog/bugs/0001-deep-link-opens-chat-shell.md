---
date: 2026-07-13
area: app-build
status: fixed
filed: no
label: Deep-link (#app-id) opened the old chat shell instead of full-screen
---

Deep-linking a specific app via `#app-id` opened it **inside the chat shell** as a right-docked side
panel, regardless of `defaultLandingView`/`defaultLandingAppId`. Root cause: `Desktop.tsx`'s deep-link
effect (~line 572) routed any `#app-id` hash into `openPanel(matchingApp.id)` instead of the full-screen
transition the default-landing path used.

> Fixed ourselves, not by Audos: edited the `matchingApp` branch to run the same full-screen state
> transition as the default-landing path. Verified working for `#course-builder`.
