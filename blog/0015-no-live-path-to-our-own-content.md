---
date: 2026-07-17
product: field-notes
status: open
label: No live path to our own content
---

# No live path to our own content

Yesterday's post ended on an update, added the same day it was written: a fix for the EmailGate
problem had shipped, marked itself `Complete`, and turned out not to be live at all. We closed that
post promising to verify by hand before calling anything done. This one is that verification,
carried one step further than we expected to have to go.

## The question we actually needed answered

The invite-link rewrite (`bugs/0027`) was already known broken — cold-loading the real URL still
showed the old email prompt, a made-up invite token did nothing, and the served source had zero
trace of anything new. Before doing anything else — writing more code, filing more bugs, minting
invite tokens that would go nowhere — there was one question that mattered more than all of them:
with the new system broken, does the *old* path still work? Can anyone, right now, actually reach
the six real content items sitting in our database?

We didn't answer that from memory or from a prior session's notes. We asked Otto to walk it cold, in
a real headless browser, step by step, and report exactly what rendered at each one — not what the
code implies should happen.

## What actually happened, cold

A never-visited browser loaded the bare live URL: the generic marketing hero, same as always.
Clicking "Start now" opened Audos's own email modal. A fake, never-used email was typed and
submitted.

> No OTP code was requested. It went straight through. And what rendered next was not our
> `site_access` username/password screen — it was a generic, unbranded AI assistant shell ("How can
> I help?", "New conversation") with one sidebar app: "Build Findings Log." That app opened to "No
> findings yet" — TOTAL 0, OPEN 0, ACTIVE BUGS 0, EXPERIMENTS 0.

No password prompt ever appeared. There was nothing to sign into.

## Ruling out the obvious wrong answer

The instinctive read is "the content got lost." A direct, read-only database check said otherwise:
`content_items` still holds all six real rows, untouched. The app a fresh visitor actually lands on,
`build_findings`, reads a different table entirely — one that's always been empty. The data is fine.
The problem is that the live entry path simply doesn't go anywhere near it.

That distinction matters. This isn't a data-loss bug. It's an access-path bug: two different broken
routes (the still-unpublished invite-link rewrite, and now this) both fail to reach content that
exists and is intact the whole time.

## Where that leaves the honest answer

We were asked directly: can we reliably say a person can access their own content on this platform
right now? The honest answer is no. Not via the new system (confirmed broken, `bugs/0027`) and not
via the old manual-login workaround either (confirmed broken here, `bugs/0028`) — both dead-end
before reaching the Findings feed, the password gate, or Field Scout.

Given that, the planned full sync — the other 67 items waiting to join the six already in
`content_items` — is on hold. There's no reason to grow a pile of content on a shelf nobody can
currently reach.

## The question this raises, worth naming plainly

Two platform-owned gate rewrites in the same day, both self-reporting success, both wrong in
different ways — one served stale code, one routes to the wrong app entirely. That's not a
coincidence worth shrugging off. It's a second data point, on the same day, for the same underlying
thing this whole project keeps finding: on a platform where the builder doesn't control the publish
step, "it says Complete" and "it works" are two separate claims, and closing the gap between them
takes actual verification, every time, no exceptions for the fix that was supposed to fix exactly
that.

It's also the moment a different question got asked out loud for the first time: if access control
for a page meant to be invite-only, hand-picked, and password-gated keeps landing on infrastructure
we don't control and can't reliably verify, does this particular piece — the findings record itself —
need to keep living here at all? That's not decided. But it's now on the table, not routed around.
