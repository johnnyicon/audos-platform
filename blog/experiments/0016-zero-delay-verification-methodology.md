---
date: 2026-07-13
area: app-build
status: corrected
label: "Zero-delay screenshots caught a shell-flash bug a delayed check had already \"verified\" clean"
---

**Hypothesis:** Was the full-screen shell fix actually complete, or would a stricter verification method
— a screenshot fired with zero delay, versus one fired a few seconds after navigation — surface a gap
the earlier "verified live" check had missed?

**Method:** After fixing `#app-id` deep-link routing and a separate "return to chat" teardown bug,
verified two ways on the same build: a delayed screenshot (the method already in use, which had reported
clean), and a screenshot with zero delay immediately after navigation.

**Result: corrected an earlier "verified" claim that turned out to be wrong.** The zero-delay method
caught a real, reproducible flash of the old chat shell on initial paint — invisible to any check that
waits even a few seconds first. Root cause: `useState` shell-mode defaults were only corrected by a
post-paint `useEffect`, a textbook effect-driven-correction flash. Fixed by moving the resolution into
lazy `useState` initializers that read `window.location.hash` synchronously, no effect needed. Re-verified
with the corrected zero-delay method across two apps, twice each. Also surfaced a genuine, still-open
per-app inconsistency — one app never actually picked up the deep-link fix — filed as its own bug rather
than folded into the "fixed" claim.

See `BACKLOG.md #8`; `docs/platform/22-eliminating-the-chat-shell-playbook.md`, verification checklist
step 5; `blog/0008-three-shell-bugs-and-a-verification-blind-spot.md`.
