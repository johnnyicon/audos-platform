---
date: 2026-07-16
area: app-build
status: fixed
filed: yes
filed_ref: "Priority Support, field-notes workspace, 2026-07-16 — escalated to a human engineer, then retracted same day after we found our own mistake"
audos_status: pending
label: "RETRACTED: what looked like file writes vanishing on a crash was actually our own verification checking the wrong endpoint"
---

**Retracted the same day it was filed.** We originally reported that a build job's file writes vanish if
it errors before some final step. They don't — we were checking whether the *published* bundle had
updated, not whether the *draft* had. The app was there the whole time.

> Original claim: a ~2-hour job wrote a full app, live-tested it working, then errored before posting a
> completion report — and `GET /api/space/{id}/files` came back `{"files": []}` immediately after,
> which we read as data loss. What we'd actually done: that endpoint only ever reflects the published
> bundle, never draft/unpublished work. Proof: the workspace's own Preview panel, checked side by side —
> "Live" mode shows the same old marketing page the API implied (never updated), while "Draft" mode shows
> the real, working app. We logged in ourselves with the test credentials and reached the actual content
> feed. A follow-up job confirmed the original job's app file, config, and hooks were already present when
> it started — it made one small fix rather than rebuilding from scratch, meaning the first job likely
> never lost anything; it crashed before *it* found the right draft-preview URL to confirm and report its
> own success, and we compounded that by checking the same wrong, publish-only endpoint independently.
> **What's still a genuine, smaller finding:** there's no obvious way to tell "nothing was built" apart
> from "something was built but only exists as an unpublished draft" without knowing to check the Preview
> panel's Draft toggle specifically — worth documenting, but it's a discoverability gap, not data loss.
> We followed up on the original ticket (already escalated to a human engineer) to correct the record
> before real engineering time got spent on a bug that wasn't real.
