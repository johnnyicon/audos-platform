---
date: 2026-07-13
priority: 1
status: not filed
label: Automated post-build smoke check before a job may say "verified"
---

Even a minimal automated check — load the affected route, confirm no thrown error, confirm the target DOM
element exists — run **before** a job is allowed to report "verified," would have caught both of this
week's false "verified" claims immediately, with no human re-test needed.

This is the single highest-leverage request in this whole list: everything else here is smaller than
closing the gap between what a job reports and what's actually live. See bug `0007` (job self-reports
don't reliably match live reality) for the incidents that motivate this.
