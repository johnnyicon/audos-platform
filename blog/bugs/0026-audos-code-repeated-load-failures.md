---
date: 2026-07-17
area: app-build
status: open
filed: no
label: "Audos Code (Beta 0.3.0) failed to load 4 times across roughly 40 minutes of normal use"
---

Across two work sessions the same day, `audos.com/portfolio/code` failed to fully load four separate
times: twice with an explicit "Audos Code could not be loaded. Try again." error (recovered on manual
retry), and twice with a "Hang in there, there's a bit of traffic. We're still connecting Audos Code."
message that never resolved on its own and required a hard navigation to recover from — one of those
two never recovered at all within the session.

Not filing this as a severe bug — it's a beta product (labeled `BETA 0.3.0` in its own UI, with an active
July 2026 "What's New" changelog), and every failure was recoverable by reloading. Logged because the
failure rate (4 in ~40 minutes of actual use, not stress-testing) is high enough to be a real, current
data point about the surface's maturity, not a one-off. Worth revisiting once the beta stabilizes —
if this same session's other findings (`experiments/0020`) are representative, the underlying capability
is genuinely valuable; the delivery is just not yet reliable.
