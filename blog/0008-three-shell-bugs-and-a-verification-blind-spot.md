---
date: 2026-07-13
product: DoKnow
status: fixed
label: Three shell bugs, one blind spot
---

# Three shell bugs, and a verification method we didn't know we needed

The full-screen fix from the day before felt done. It wasn't — not because the fix was wrong, but
because "verify live" turned out to have more than one meaning, and we only knew one of them.

**Bug one: the fix didn't apply the same way to every app.** `#course-builder` rendered clean and
full-screen every time. `#doknow-mockup-test` — same code path, same config shape as far as we could
tell — kept opening inside the old chat shell, even after a fresh cache-bust ruled out stale caching.
Three separate diagnostic jobs tried to diff the two apps' config entries to find the difference. All
three came back truncated at the exact point the answer would have appeared. We filed it as a real bug
with Audos rather than keep guessing, and since we no longer needed that app anyway, we deleted it —
the underlying inconsistency stays open for whoever builds a second app on this platform next.

**Bug two, and a mistake worth admitting.** The small "Assistant" pill — meant to keep chat reachable
without it dominating the screen — reappeared as the full old three-pane interface when clicked. Shown
a cropped screenshot first, we called this fine: "just a preserved affordance, working as intended."
Wrong. A fuller screenshot from the other side of the table showed the truth plainly — a complete
navigation back to the old chat+dock+panel layout, not a small overlay. Correcting our own read in the
moment mattered more than being right the first time. The actual fix, once diagnosed properly: the
button was calling a state-teardown function that dismantled the whole full-screen shell; changed to
open a small popup instead, and this time verified by actually clicking it before calling it done.

**Bug three was the one that taught us something durable about our own process.** Even after both fixes
above, on a genuinely fresh page load, the old dock would flash into view for a fraction of a second
before disappearing. Every one of our "verify live" checks up to that point waited a few seconds before
looking — which is exactly long enough to miss a flash that resolves in milliseconds. Caught only when
someone watched a real, un-delayed page load with their own eyes. We reproduced it deliberately —
screenshotting immediately after navigation, zero artificial delay — and there it was: the dock,
rendered, gone a moment later.

The cause was a clean, unglamorous React pattern once we actually looked: the relevant state defaulted
to shell-mode on mount, and a `useEffect` only corrected it to full-screen *after* the first paint had
already happened. Classic flash-of-wrong-content. The fix moved that decision into a lazy state
initializer that reads the URL synchronously, before the first render — no effect, no correction pass,
no flash. Verified this time with the method that actually catches it: screenshot immediately after
navigation, not after a polite pause, across two different apps, twice each.

**One more thing worth recording plainly, since it came out of testing all of this.** We ran an
experiment pinning a specific model (`fable-5`) to a build job instead of leaving it on auto-pick, and
had the resulting app report what it believed itself to be. It came back correctly hedged: *"believes it
is 'Fable 5'... per its system prompt and internal model slugs like 'claude-fable-5-thinking' — this is
a system-prompt-provided identity, not independently verifiable by the agent itself."* Model pinning is
real and works. And the honest self-qualification in that answer was itself useful evidence that Audos
genuinely injects system-level context into these jobs, not just a raw forwarded prompt.

The throughline across all three bugs: every one of them survived at least one round of "it's fixed"
before it actually was. Not because anyone was careless — because "verified" kept meaning something
narrower than we assumed, one layer at a time.
