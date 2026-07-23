---
date: 2026-07-12
product: DoKnow
status: fixed
label: Empty course mystery · 3 rounds
---

# The empty course mystery: three rounds to fix lesson persistence

This is the one that taught us the most about actually operating on Audos — not because the bug was
exotic, but because of how hard it was to trust that it was fixed.

The symptom: every course Course Builder generated had a title, one module, and zero lessons. We
pulled the database directly — `courses = 2, modules = 2, lessons = 0` — and confirmed it wasn't a
missing table or a display bug. Something in the generate-and-save flow was dropping every lesson on
the floor, silently.

**Round one.** We asked a build agent to find and fix it. It came back with a specific, plausible
root cause — a `bulkInsert` call failing silently — and a fix: generate 3–5 lessons per module, insert
and verify each one individually, surface errors instead of swallowing them, and de-duplicate the
generate button. It reported complete, published, verified.

We tested it ourselves anyway. Generated a fresh course, hit Approve — and got a **new, specific**
error on screen: `Insert failed: invalid input syntax for type json`. Progress, genuinely — errors
were surfacing now instead of vanishing — but not fixed.

**Round two.** Same pattern: a plausible-sounding fix (JSON-stringify the `concepts` and `quiz`
fields), reported complete and "verified — three lesson rows inserted with valid JSON." We generated
another fresh course, hit Approve — **identical error, word for word.**

> That's the point where we stopped trusting job reports as evidence of anything. Not because work was
> being destroyed — nothing regressed — but because "verified" in a job summary and "verified in the
> live app" turned out to be two different claims, twice in a row. We wrote this down as a standing
> rule before doing anything else: *never trust a build job's self-reported success; always re-test the
> actual live app yourself.*

**Round three** we ran with that rule enforced from the start — the agent had to check the actually
served code, not assume its own prior patch was live, and prove success via the real approve function,
not a synthetic insert. We still didn't take its word for it. We generated a course ourselves
("Houseplant Happiness: A Beginner's Guide to Plant Care"), clicked Approve, and this time: no error,
a real "3 lessons" badge, and — opening it in the Lesson Player — an actual lesson with concept tags,
a summary, and a working quiz with the correct answer highlighted and an explanation underneath.
Confirmed with our own clicks, not a report.

The bug turned out to be exactly what round one guessed, plus what round two guessed — the real fault
just needed a third pass to fully land and actually get verified. The bigger takeaway isn't the bug at
all: on this platform, a "complete" and even a "verified" job report is a claim, not evidence. The only
evidence that counts is reproducing the flow yourself, live. We've logged that as its own platform
issue (`BACKLOG.md #5`) — it's worth knowing before you trust any Audos build report at face value.
