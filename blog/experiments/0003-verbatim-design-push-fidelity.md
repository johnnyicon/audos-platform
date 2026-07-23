---
date: 2026-07-12
area: app-build
status: confirmed
label: Does a hand-built mockup survive a "port verbatim, don't redesign" build job?
---

**Hypothesis:** A build agent handed a complete, self-contained HTML/CSS/JS mockup and explicitly told
to port it verbatim will actually preserve it, rather than reinterpreting a design brief in its own
style.

**Method:** Sent a 903-line self-contained mockup to a build agent, with the "port verbatim, don't
redesign" instruction repeated twice for emphasis. Verified by opening the live app directly rather than
trusting the job's own completion report.

**Result: confirmed pass.** Matched the mockup element-for-element — wordmark, streak pill, nav rail,
lesson card, all of it. Left open at the time: whether a verbatim port can also be wired to real database
data in the same pass, or whether design fidelity and data wiring trade off against each other. Answered
later by the "Red Pill" migration experiment.

See `blog/0007-design-fidelity-can-we-push-our-own-ui.md`.
