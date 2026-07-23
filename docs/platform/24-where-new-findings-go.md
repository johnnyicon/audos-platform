# Where a new finding actually goes

*Written 2026-07-13, extended 2026-07-15 with the archive/changelog workflow, corrected same day once*
*the activity log itself moved out of this repo. The decision rule this SDK has been following ad hoc,*
*made explicit so the next agent doesn't have to guess.*

> **This repo does not hold the activity log.** `ACTIVITY-LOG.md` lives in `doknow-kb/audos/`, not here —
> this repo (`audos-platform`) is meant to stay a clean, distilled, shareable SDK record; the raw,
> product-specific work-in-progress log doesn't belong in it. See `doknow-kb/audos/README.md` for the
> full reasoning, and the `audos-sdk-*` skills (which live in `doknow-kb/.claude/skills/`) for the actual
> tooling — they operate cross-repo: source in `doknow-kb`, destination here.

Five places a piece of information can land. Pick by what kind of thing it is:

| It's a... | Goes in | Example |
|---|---|---|
| **Fast-scan record of an action taken** (dispatched a job, got a result) | `doknow-kb/audos/ACTIVITY-LOG.md` (a different repo — see note above) | "Dispatched X, model Y, job #Z, result: ..." — one row |
| **A tracked bug/issue** with a fix state to follow over time | `BACKLOG.md` (+ a file in `docs/platform/bug-reports/` if formally filed with Audos) | Deep-link routing inconsistency |
| **Durable platform knowledge** — a mechanism, a pattern, a playbook someone will want to *find* later, not just skim once | `docs/platform/NN-topic.md`, next available number, one row added to `skill/SKILL.md`'s index table | The dispatch-then-poll SOP, the chat-shell playbook |
| **A narrative worth telling** — the story of how something was found/fixed, for the cohort/Audos to read | `blog/NNNN-slug.md`, per `blog/HOW-TO-UPDATE.md` | Any blog entry |
| **A terse "what changed in the SDK's own guidance" summary**, for someone who wants the delta without the story | `CHANGELOG.md`, newest entry on top | "Corrected docs/platform/06's storage endpoint claims" |

**A single finding often touches more than one of these** — that's normal, not a sign you picked wrong.
Example: the flash-of-old-shell bug got a **log line** (what was done), a **BACKLOG entry** (tracked fix
state), a **durable doc section** (the React mechanism + the fix pattern, reusable by anyone), and a
**blog entry** (the story, including the part where we got it wrong once first). Write to all of the
relevant ones, not just one.

**When in doubt:** if you'd want to *grep for it later*, it belongs in `docs/platform/`. If you'd want to
*read it as a story*, it belongs in the blog. If you just need to know *what happened, when*, the
activity log. If someone needs to *track whether it's fixed*, the backlog. If you want the terse "what's
different now" delta without re-reading the story, the changelog.

## The harvest pass — the standing loop this SDK runs periodically

**Run this via the `audos-sdk-harvest` skill** (in `doknow-kb/.claude/skills/`) rather than by hand —
it encodes all four steps below plus the exact file formats. This section documents what it does, for
reference.

`doknow-kb/audos/ACTIVITY-LOG.md` is meant to be trimmed regularly, not left to grow forever — its own
header says it's the "fast-scan" record, and a log that never gets pruned stops being fast to scan.
Every harvest pass does, in order:

1. **Review the unharvested rows** in the live `doknow-kb/audos/ACTIVITY-LOG.md`. For each, promote
   whatever's durable — and genuinely generic to the platform, not DoKnow-specific business detail —
   into this repo's `docs/platform/`, `BACKLOG.md`, or a bug-report file, per the table above.
   DoKnow-specific process detail with no generic platform lesson in it stays out of this repo entirely.
2. **Write a blog entry** narrating the pass itself — not a duplicate of the durable docs, the story of
   what the sweep found (what got reclassified, what was already captured, what gaps turned up). Skip
   this step only if the pass turned up nothing narrative-worthy.
3. **Write a `CHANGELOG.md` entry** — newest on top, dated, listing exactly what changed (which docs got
   corrected, which new findings got added, which backlog items got filed), with a one-line pointer to
   the blog entry if one was written.
4. **Archive the harvested rows**, back in `doknow-kb`. Copy the live `audos/ACTIVITY-LOG.md` to
   `audos/archive/ACTIVITY-LOG-<YYYY-MM-DDTHHMMZ>.md` (full UTC date+time, not just date — a harvest pass
   can happen more than once a day), add a one-line header noting it's archived and pointing at the
   matching `CHANGELOG.md` entry (in this repo), then reset the live log to just its header/standing-rules
   with an empty table. Nothing is deleted — the archive file keeps the full raw chronological record;
   the live log just stops carrying already-harvested weight.
