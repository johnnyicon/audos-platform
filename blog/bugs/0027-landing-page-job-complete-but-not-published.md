---
date: 2026-07-17
area: app-build
status: open
filed: no
label: "A landing-page rewrite job reported \"Complete\" — the live site still served the old code"
---

Dispatched a job to rewrite `EmailGate.tsx` (the platform-owned signed-out landing view of a unified
space — see `bugs/0024`). The job finished, marked itself `Complete`, and its own outcome text claimed
specifics: the target file was the only one changed, an internal version marker was preserved, and a
"smoke-test query" had confirmed the new logic worked against the database.

None of that matched the live site.

> Three independent live checks, not taken on the job's word: (1) loading the bare production URL cold,
> with no cookies, still rendered the old Audos email-capture prompt — not the new invite-only notice.
> (2) Loading with a deliberately bogus invite token rendered the identical old page — no new branch
> exists to even reject the bad token. (3) Reading the actual served page source directly (not the
> job's diff, the real bytes a visitor's browser receives) confirmed it: the old version marker was
> still present, the old email-gate button text was still present, and a text search for any string
> from the new code — `invite`, `invite_token`, `Sign in manually` — returned zero matches anywhere in
> what's actually being served.

Root cause not fully confirmed, but a real lead: a second job had been dispatched at an overlapping time
against the *same file*, and the server logs showed the live site had been "restored... from GCS" from an
already-published snapshot — consistent with the edit landing in source but the publish/recompile step
either never running or being overwritten by the overlapping job. The database half of the same change
(a new column, dispatched as a separate job) *did* land correctly and was independently confirmed live —
so this isn't "the platform can't make database changes stick," it's specifically the file/publish path
for this landing-page component.

This is the same failure class this whole project keeps finding — a job's own completion report is not
evidence — reproduced here on our own fix for exactly that problem, which is its own small irony worth
keeping in the record rather than smoothing over.
