---
date: 2026-07-14
product: DoKnow
status: pass
label: Taking the red pill
---

# Taking the red pill

Entry five, "Escaping the shell," proved Audos *could* render a real product full-screen, no chat wrapper.
Entry nine went further and consolidated everything into one fully-owned app. Both were real progress, but
both did the same thing under the hood: build a new app, the right way, and retire the old one. Neither
answered the harder question — can an app the platform already generated on its own terms be pulled out of
the matrix *in place*, keeping its id, its file, its data, without starting over? Blue pill is staying in
the simulation without knowing it. Red pill is the one that gets you out. This entry is about actually
taking it, on an app that was already in.

What "prove it" actually meant here, stated plainly rather than left implicit: four checks, on the same
app, in this order. **Mechanics** — does it load full-screen, no dock, no flash of the old chat shell, on
a true zero-delay screenshot. **Design** — does it show our actual product design, not Audos's own
generated styling. **Data continuity** — does it still read the same real, shared data once its UI is
replaced, rather than resetting to sample content. **Scope** — does touching this one app leave every
other app in the workspace untouched, proving this was a contained, single-app operation. All four had to
pass, independently verified live, not trusted from a job's own completion report, for this to count as
proven rather than claimed.

Worth being precise about vocabulary here, because the finding below turns on it: a **workspace** is the
one top-level container for a whole business (there has been exactly one DoKnow workspace this entire
build — nothing here ever spun up a second one); an **app** is a smaller thing living inside it, one of
possibly many. Every "app" named in this build log — `doknow-app`, `doknow-home`, `one-good-thing`, the
now-deleted `course-builder` and friends — is a sub-unit of that single workspace, not a workspace of its
own.

The app we meant to use for this was the very first one this whole project started around. It had already
been deleted, in an earlier cleanup we'd asked for ourselves. So we recreated a plain stand-in the same
way the original was built — `doknow-home`, describe-then-generate, no special instructions, meant to be
an authentic "born in the matrix" baseline.

It refused to be broken. A zero-delay screenshot — the same test that caught a real flash-of-old-shell bug
earlier this week — showed it rendering full-screen, no dock, no flash, from the very first frame. That
wasn't the app being unusually well-built. It's that `Desktop.tsx`, where all four of the earlier shell
fixes live, is shared across the whole workspace, not owned per app. Fix it once and every app built after
inherits it, permanently, without being told. That's a good outcome dressed up as a failed experiment: the
mechanics escape isn't a per-app chore anymore in this workspace. It already happened, once, for good.

What hadn't escaped was the design. The baseline app picked up the workspace's teal brand color
automatically, which made it look closer to ours than it was — a different greeting, a different card
layout, a different everything else. Generic, just tinted to match. That's the real gap left to close, and
closing it in place — not in a parallel app — was the actual point of this entry.

We briefed a build agent to port our verbatim mockup into that exact file, replacing its UI, explicitly
told to preserve the real-data wiring already there rather than reset it to sample content. First pass
looked done on the report, but a check of what had actually been embedded in the brief showed only the CSS
was pasted in literally — the HTML body and script were "referenced," meaning the agent would have been
reconstructing markup from a description, not porting it. That defeats the whole technique. We hadn't
caught this from the report; we caught it from checking our own brief. Sent a corrected version with the
complete CSS, HTML, and script pasted in whole, pinned `fable-5`, and let it queue in behind the first
attempt so its write would be the one that stuck.

The completion report did what these reports always do — cut off mid-sentence at the exact line that would
have confirmed real data wiring, which model ran it, and publish status. By now that's not surprising, it's
routine. We didn't wait on it. Loaded the app cold: zero-delay screenshot, clean, no shell flash, and this
time our actual design — left icon rail, teal mark, streak pill, the real greeting copy. Let it settle, and
the data underneath was the same real, shared state the other app already had: one-day streak, the
beekeeping course at 33%, the same six courses, the same up-next lesson. Checked the app we didn't touch —
still the default, still showing the identical real data, unaffected. One small mismatch turned up between
the two apps' "cards due for review" counts, worth chasing later, not the kind of thing that changes the
conclusion.

The conclusion: an app the platform generated on its own terms, chat shell and all, can be pulled fully out
in place. Not rebuilt next to itself and switched over — the same file, the same id, migrated. Mechanics
turned out to be a one-time, workspace-wide fix, already done. Design still needs the verbatim push, every
time, per app — but it can be pushed onto an app that already exists, carrying its real data forward with
it. That's the whole shape of the escape, closed on both ends now: build clean from the start, or take an
existing app and get it out. Either way, nothing about the matrix is load-bearing. It was never the only
way through.

All four checks from the top of this entry passed: full-screen with no shell flash, our actual design in
place of the generic one, the same real data carried forward rather than reset, and the untouched app left
exactly as it was. That's what "proven" means here — not claimed, checked by hand, on the live app.

One honest loose end this surfaced: entry 10 ("Building it right the first time") had claimed a clean,
from-scratch app avoided four separate shell/data bugs purely by being briefed correctly upfront. This
entry's finding — that shell mechanics are inherited at the *workspace* level, not tested fresh by every
new *app* — means three of those four were never actually isolated as a test; they'd already been fixed
in this workspace before that app was ever built. Entry 10 has been corrected in place rather than left
to stand. The one part of it that was a real, unconfounded test — a database-write convention avoiding
its own bug class when briefed upfront — still holds. What's still genuinely untested, and deliberately
left for later rather than rushed: whether a truly fresh workspace, one that has never had its own
`Desktop.tsx` touched, can avoid the shell bugs the same way purely by being briefed correctly from
creation. That's a different experiment than anything run so far.
