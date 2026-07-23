---
date: 2026-07-13
product: DoKnow
status: pass
label: Escaping the matrix
---

# Escaping the matrix: one app, no remnants, real data

Every fix in the previous entry was a retrofit — take an app the platform generated on its own terms,
and correct it after the fact. That works, but it means fighting the same category of problem on every
app, one at a time, forever. At some point the better move stops being "patch what Audos built" and
starts being "build the thing ourselves and use Audos as a deployment target, not a design authority."

So we did that. One new app, `doknow-app`, built entirely from a technique we'd already proven works:
port our own complete, already-designed HTML/CSS/JS verbatim, with an explicit instruction not to
redesign, restyle, or reinterpret anything. That much we'd done before, for a throwaway test. This time
we took it further — wired the same verbatim UI to real data instead of static sample content, and
folded three separate generated apps' worth of functionality into one app we actually own.

Concretely: the home shelf now reads real rows from the same `courses`/`modules`/`lessons`/
`learner_progress` tables the old apps used — not fixture content. The course-generation flow — name a
topic, get a leveled course, review an editable preview, approve it — was ported in rather than left in
its own separate app. So was the lesson player: real content, a real scored quiz, and a mark-complete
action that actually writes to the database.

We tested all three flows by hand rather than take a completion report's word for it, because by this
point in the week that would have been a mistake on principle. Generated a real course — "Buzzing
Basics: An Introduction to Urban Beekeeping" — approved it with no error (this is exactly where the
platform's course-generation bug lived originally; it's clean now). Took the quiz, got a real score with
real explanations. Clicked Mark Complete, and hit the same `Insert failed: invalid input syntax for type
json` error we'd already fixed once in a different app — a fresh code path, same unserialized-JSON
mistake, third time today. Fixed the same way as before (find the offending column, `JSON.stringify` it),
re-tested by hand, and watched it actually work: streak ticked from zero to one day, the course's
progress moved from 0% to 33%, and the "up next" card correctly advanced to the second lesson.

Then we made it the default. The bare workspace URL now boots straight into this one app — not a
chat-first shell, not a scaffold sitting next to the real thing, the actual product.

Two things we're carrying forward as standing practice rather than one-off choices for this build. Every
job dispatched to the platform now pins a specific model explicitly (`fable-5`, with a named fallback) —
never auto-pick, and every dispatch gets logged with its model, its task ID, and its outcome in a
running activity log kept separately from this narrative record. And every job brief now states its
instructions in full, every time, rather than assuming any repo convention file will be picked up — a
discipline that turned out to be load-bearing for reasons entirely outside our control (see the AGENTS.md
loading bug noted in the Cursor backend research).

The apps this one replaces — the originally generated Course Builder, Lesson Player, Coach Queue,
DoKnow Home, and the two throwaway experiments — have been deleted, and we checked: the workspace now
has exactly one real app, and a deleted route like `#course-builder` falls through gracefully to the
one that's left rather than breaking. What's left of the platform's own defaults in this workspace,
after today, is the database layer underneath — which was never really the problem. The problem was
always the shell.
