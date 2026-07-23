---
date: 2026-07-17
area: app-build
status: open
filed: no
label: "A cold visitor who clears the email gate never reaches the password screen or the real content"
---

Walked the full real-world path end to end in a cold headless browser (no cookies, no prior
session) to answer one question: can a person reliably reach field-notes' actual content right
now, via any live path — not just the still-broken invite-link rewrite (`bugs/0027`).

> (1) Bare live URL, cold: the generic marketing hero (`bugs/0024`'s EmailGate) loads. (2) Click
> "Start now": the Audos email modal appears. (3) Type a fake, never-seen email and submit: it goes
> **straight through — no OTP requested**, and the email step creates a new CRM contact and fires
> funnel/ad-tracking events (same undisclosed side effect as `bugs/0024`). (4) What renders next is
> **not** the site's own `site_access` username/password gate — it's a generic, unbranded AI
> assistant shell ("How can I help?", "New conversation") with a single sidebar app, "Build Findings
> Log." (5) That app shows "No findings yet" — TOTAL 0, OPEN 0, ACTIVE BUGS 0, EXPERIMENTS 0. No
> password prompt ever appeared to test against.

A direct, read-only database check ruled out data loss as the cause: `content_items` still has all
6 real rows intact. The app a fresh visitor actually lands on, `build_findings`, reads a different,
empty table. So this isn't "the content disappeared" — it's that the only entry path currently live
for a cold visitor doesn't lead to the Findings feed, `site_access` gate, or Field Scout at all; it
dead-ends at an unrelated, empty app.

Net effect: right now there is no live path — not the invite-link rewrite (broken per `bugs/0027`),
and not the "old" email-then-password workaround (this bug) — by which an outside visitor can reach
field-notes' real content. The `site_access` password gate and the Findings feed may still exist
in the signed-in space's source, but this session could not confirm it — there was no screen to
reach them from. Root cause not yet established: could be the same
publish/snapshot problem as `bugs/0027`, or the signed-in space's routing/gate wiring genuinely
changed. Logged as its own finding rather than assumed to be the same bug, since the failure mode
(dead-ends at a *different* app, not "shows old code") is distinct from 0027's.
