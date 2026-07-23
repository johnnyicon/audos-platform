---
date: 2026-07-13
area: app-build
status: fixed
filed: no
label: Flash of the old chat shell on initial page load
---

Even after the deep-link and Assistant-pill fixes, every fresh page load **visibly flashed the old
dock/sidebar** before the corrected full-screen view took over. Invisible to a normal check (screenshot
a few seconds after navigating looks clean) — only caught by screenshotting with **zero delay**,
immediately on navigation, a method a user's direct challenge forced us to adopt.

Root cause: `Desktop.tsx`'s `useState` calls default to shell-mode (`isSidebarOpen=true`,
`activePanelId=null`, `mobileView='chat'`), and the deep-link `useEffect` only corrects this **after**
React's first paint — a classic effect-driven-correction flash.

> Fixed ourselves: replaced the plain `useState(default)` calls with lazy initializers that resolve from
> `window.location.hash` synchronously, so the correct state is set before the first render — no effect,
> no correction pass, no flash. Verified with the zero-delay method across two apps, twice each.
