---
date: 2026-07-12
area: app-build
status: confirmed
label: Can our real frontend stack run inside an Audos app at all?
---

**Hypothesis:** Audos compiles a single file per app via server-side ESBuild and resolves dependencies
via CDN, not `npm install`. Does the team's actual stack — TanStack Query, GSAP, Radix, three.js,
react-three-fiber — survive that pipeline, or does the platform's build model quietly rule some of it
out?

**Method:** Built a disposable "Lib Probe" app that live-imports and renders each library with a
pass/fail badge and the real error text if one failed.

**Result: mostly confirmed pass.** 4 of 5 rendered cleanly — react-query, GSAP, Radix, and raw three.js
all worked. react-three-fiber was flaky on cold start (hit a watchdog timeout once) but worked on
reload, so it's a qualified pass rather than a clean one: usable, but worth remembering it can be slow
to warm up.

See `blog/0005-escaping-the-shell.md`, "Experiment 1: the library probe."
