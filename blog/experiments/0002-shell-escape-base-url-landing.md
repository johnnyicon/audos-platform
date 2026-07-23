---
date: 2026-07-12
area: app-build
status: confirmed
label: Can a full-screen app become the workspace's base-URL landing experience?
---

**Hypothesis:** Setting `desktop.layout.defaultLandingView`/`defaultLandingAppId` makes a single app the
workspace's base URL, with zero chat shell visible — not just reachable from inside the shell.

**Method:** Built a standalone "DoKnow Home" app, set the config, published, then loaded the bare
workspace URL cold — no deep-link hash, no prior session.

**Result: confirmed pass**, with one caveat filed as its own open item: this worked for the base URL
only. Deep-linking a specific app via `#app-id` hash still opened it inside the chat shell as a
right-docked side panel — the two paths aren't the same code, and only one of them was fixed by this
config. Became `BACKLOG.md #8`.

See `blog/0005-escaping-the-shell.md`, "Experiment 2: escape the shell."
