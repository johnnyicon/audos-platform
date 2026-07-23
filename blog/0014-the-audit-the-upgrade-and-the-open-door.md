---
date: 2026-07-17
product: field-notes
status: pass
label: The audit, the upgrade, and the open door
---

# The audit, the upgrade, and the open door

Yesterday ended with a question nobody could quite answer cleanly: after two weeks of building real
things on Audos — a working product in DoKnow, a stress-test in field-notes, a growing pile of bugs and
workarounds — is this platform actually architected for the kind of development a coding agent does, or
have we just gotten good at routing around what it can't do? The honest feeling by the end of that day
was somewhere close to "we keep making lemonade out of lemons, and that's not the same as the lemons not
being there." So today's first job was to stop routing around the question and just answer it.

We asked a fresh session — no memory of any of this, reading only the corpus itself — to do a red-team
audit: read every bug, every experiment, every correction we'd logged, and say plainly whether Audos is
built for hardcore agentic development. Not counting bugs. Finding the small number of root causes behind
all of them. And because this was headed to one of Audos's own co-founders, not staying internal, it had
to survive an adversarial read before it went anywhere — every claim checked against the evidence, every
absolute softened to what the evidence actually supported, nothing sent that a five-minute look at the
real product would embarrass us on.

That last part turned out to matter more than expected. Checking the audit's own claims against the live
platform, one thing didn't add up: a whole second product surface, "Audos Code," versioned and actively
changing, that the corpus had never once mentioned. Not a rumor — a real, running thing, sitting one
click away in the same workspace we'd been building in for two weeks, and nobody had opened it.

## What we found when we actually opened it

Two things, in tension with each other, which is usually how the most useful findings arrive.

The good one first: Audos Code streams. Turn on the right setting and you watch real, incremental,
token-by-token output while an edit runs — not a spinner, not a poll, an actual live log of what's
happening. And the backend underneath it isn't the Cursor pipeline every other finding in this project
has been written against — it's Claude Code, Anthropic's own agent, with a routing panel that shows
which model picked up which turn. That's a genuinely different architecture than the one we'd spent two
weeks characterizing, and the audit needed to say so plainly rather than quietly keep citing Cursor-shaped
failure modes against a system that had already moved past them.

Then the other one. To test whether the new surface actually behaved differently, not just looked
different, we gave it the smallest possible edit — change a period to an exclamation point in a real
headline — and watched it work. Streamed activity, real tool calls, a confident close: *"Done. The
headline now reads 'Stop saving. Start knowing!'"* We zoomed into the actual page. Still a period.
Refreshed by hand. Still a period. The exact failure this whole log keeps finding — a job that reports
success while nothing changed — had just happened again, live, in the surface built to fix exactly that
problem. Streaming solved observability. It didn't touch trust. Those turned out to be two different
things wearing one name.

## Turning the discovery into use

With the audit's evidence base stronger for having actually looked, there was a second, more practical
question sitting right behind it: `content_items` had been empty since the day it was created, blocked
by a direct-Postgres credential flow that was 0-for-3 across every attempt anyone had made. Could this new
surface — the one that had just lied to us about a punctuation mark — also write real data reliably?

First attempt said no in an interesting way. Sent six real rows as ordinary multi-paragraph prose; the
agent came back and said, plainly, that it hadn't actually received the row data — and rather than
inventing something to fill the gap, it went and checked the workspace and the database itself for a
seed source before asking us to resend. That's the discipline this whole project has been asking Audos
to build in one form or another for two weeks, showing up unprompted, in the same session that had just
failed to move a punctuation mark. Sent the same six rows again as compact JSON instead of prose. Landed
cleanly. Row count went from zero to six, and the agent queried the table back itself before calling it
done, rather than just asserting a number.

We didn't take that as the final word either. Checked the workspace's own preview panel by hand — all
six items rendered, correctly split into the right sections with the right status. Then checked the
actual live site, no publish step taken, and asked its own embedded assistant to show the latest bugs. It
answered correctly, using content that had existed for all of about ninety seconds. Database writes reach
Live the instant they happen; only the file bundle needs a publish. The credential bug blocks one door.
This is a second one, and it works.

## The door nobody built

Then a real visitor — one of us, using our own account — tried the live site and hit something neither
of us had put there: an email prompt, asking for an address we already knew, before anything we'd
actually built ever showed up.

Rather than guess, we asked Otto to go read the served code directly and report back literally what it
found, not what seemed plausible. What came back was more consequential than expected. Every workspace on
this platform has exactly one signed-out view, and it isn't yours — it's Audos's own `EmailGate`
component, unconditionally the front door, with nothing else in the workspace able to mount until it says
it's done. The password screen we'd built, the one this whole site existed to put in front of visitors
instead of relying on someone else's sharing model, had never actually been the first thing anyone saw.
It was the *second* gate, structurally incapable of being anything else, sitting behind a screen we
didn't build and hadn't audited.

And that screen isn't neutral. Submitting an email into it registers a CRM contact and fires ad-tracking
pixels to Meta and Reddit — because it's Audos's own lead-capture funnel, quietly doing lead-capture
things, dressed as a generic sign-in step. A private tool built specifically so we controlled who saw it
and what happened to their information had been handing that information to two ad platforms without
either of us knowing, the whole time it existed.

We're fixing it the way that actually matches what this site was supposed to be from the start: killing
that gate entirely and replacing it with a closed-alpha system built around personal invite links —
approve someone, hand them one URL, they're in, no email typed anywhere, nothing shared with anyone we
didn't choose to share it with. The build is running as this goes up; we'll verify it by hand, the same
way we verified everything else today, before calling it done.

> **Update, later the same day.** The fix ran, marked itself complete, and reported a specific, confident
> outcome — the right file changed, a version marker preserved, a smoke test passed. We checked anyway.
> Cold-loaded the real URL: the old email prompt was still there. Tried a made-up invite link: same old
> page, no new logic at all. Read the actual bytes being served, not the job's diff: the old code, word
> for word, with zero trace of anything we'd just asked for. The fix for "don't trust a job's report" had
> itself just given us a report we couldn't trust — on a Tuesday, in the same file, in the time it took to
> write this sentence. Likely cause: two jobs edited the same file at once, and the live site kept serving
> an older published snapshot underneath both of them. Logged as its own finding rather than quietly
> fixed and forgotten, because it's too on-the-nose not to.

## Where the day actually leaves things

The audit is real, better-evidenced than it would have been without today's detour, and honest about
both directions — the runtime does more than its own docs claim, and the loop around it still asks you to
distrust every success message, including, as of today, the newer ones. field-notes has actual content on
it for the first time, reachable on the real URL, verified by hand rather than taken on a job's word. And
the site's own front door, the one thing we'd assumed was simply ours because we built it, turned out not
to be — which is its own small, sharp instance of everything this whole log has been circling for two
weeks: on this platform, "we built it" and "we control it" are not the same claim, and the gap between
them is usually invisible until someone actually walks through the door.
