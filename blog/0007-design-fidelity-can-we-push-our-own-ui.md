---
date: 2026-07-12
product: DoKnow
status: pass
label: Design fidelity
---

# Design fidelity: can we push our own UI and have it stick?

Once we could generate real content and land on a real dashboard, one design question was still open.
Every app we'd had Audos build — even the successful ones — was styled by *its* interpretation of a
text brief. Reasonable results, but never quite what we'd actually drawn. So we asked the sharper
version of the question: instead of describing a UI and hoping the build agent's taste matches ours,
can we hand Audos a UI we already built, exactly as we want it, and have it survive intact?

We had a self-contained mockup on hand — a single 903-line HTML file with its own CSS and vanilla JS,
built earlier in the day to pin down what DoKnow's home screen and lesson player should actually look
like. We sent the whole file to a build agent with one instruction, stated plainly and repeated twice:
port this **verbatim**. Do not redesign, restyle, rewrite copy, or improve anything. Adapt only what
the platform's component format strictly requires.

> We did not take the completion report at face value — by this point in the day that would have been a
> mistake on principle. We opened the live app ourselves. It matched. The wordmark, the streak pill, the
> four-icon navigation rail, the greeting copy, the hero lesson card with its format chips and source
> citation, the course shelf with its progress values, the full lesson body text on a second view — all
> of it was our content, rendered as we'd designed it, not Audos's own interpretation. In-app navigation
> between views worked with no page reload.

The adaptation required to make it fit Audos's app format was genuinely minimal: wrap the markup in a
default-exported component, move the vanilla JS behavior into a React effect, inject our original CSS
byte-for-byte. Boilerplate, not a rewrite.

This changes how we'll actually work with the platform going forward. The default mode —
describe-then-generate — reliably produces the "every Audos app looks the same" effect, because the
build agent is filling in your intent with its own design judgment. Design-then-port sidesteps that
completely: the fidelity is ours, because we're not asking Audos to invent anything. If you care about
a specific look and feel, build it once, wherever you're most comfortable, and hand it over — don't
describe it and hope.

One thing we haven't tested yet: whether a verbatim-ported UI can also be wired live to real database
data in a single pass, the way our full-screen dashboard experiment was. Nothing we saw today suggests
it can't — that's simply the next thing on the list.
