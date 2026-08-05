# Prompt for a fresh Fable 5 session — Audos architectural red-team audit

> This file is the prepared prompt itself (everything below the next line), meant to be pasted verbatim
> into a new Fable 5 session that has no prior context on this project. Saved here so it isn't lost and
> so the source-material paths stay accurate as the repo evolves.

---

## Task

You're doing a red-team architectural audit of the Audos no-code platform, based entirely on evidence
already collected over roughly two weeks of real, hands-on building. This is a **synthesis task, not a
bug hunt** — every individual bug, experiment, and finding already exists in the repo below. Your job is
to read across all of it and answer one question honestly: **is Audos architecturally capable of
supporting real, hardcore agentic app development — the way building with a terminal-attached coding
agent (Claude Code, Codex) is — or does it require constant, structural workarounds that a properly
architected agentic dev platform wouldn't need in the first place?**

Don't answer that question by counting bugs. Answer it by finding the small number of *root
architectural gaps* that keep producing the same shapes of pain across many unrelated findings — the
80/20: the handful of fixes that would unlock the most real capability, versus the long tail of
individual bugs that are just symptoms of those same root gaps.

## Context

Two products have been built on Audos so far:
- **DoKnow** — the actual product, a spaced-repetition learning app. Live and published on Audos.
- **field-notes** — a small internal tool, built deliberately in its *own separate* Audos workspace
  specifically to stress-test the platform's rougher edges without risking DoKnow's workspace. Most of
  the sharpest recent findings came from this deliberate stress-testing.

Everything either project has learned about the platform — capability findings, bugs, experiments,
feature requests — has been logged in one shared repo, `audos-platform`, precisely so it generalizes
across every product built on Audos, not just one.

## Source material — read all of it before writing anything

All paths below are relative to `~/Workspace/audos-platform` unless stated otherwise.

- `docs/platform/*.md` — the accumulated capability reference and narrative findings docs (numbered,
  roughly chronological — read `19` through `25` especially closely, they cover the platform's
  operational model: job dispatch, verification blind spots, escalation paths, the chat-shell
  architecture).
- `blog/*.md` (`0001` through `0013`) — narrative write-ups of each major finding, in the order they
  happened, including corrections we had to make to our own earlier conclusions.
- `blog/experiments/*.md` (`0001` through `0019`) — every deliberate, hypothesis-driven test run against
  the platform: hypothesis, method, result, whether it was later confirmed, corrected, or left
  inconclusive. This is the highest-signal source for your audit — it's already structured as
  falsifiable claims with evidence.
- `blog/bugs/*.md` (`0001` through `0023`) — every confirmed platform bug, whether fixed, still open,
  filed with Audos or not.
- `blog/feature-requests/*.md` (`0001` through `0009`) — things we concluded don't exist and should.
- `BACKLOG.md` — the terse issue-by-issue index, useful for a fast second pass.
- `CHANGELOG.md` — dated summary of what changed in our own understanding over time; useful for seeing
  which conclusions got revised and why.
- `field-notes/ACTIVITY-LOG.md` — the raw, blow-by-blow log of the field-notes build, including several
  moments where we caught ourselves trusting an unverified claim (a job's self-report, our own earlier
  correction) and had to re-verify a second time.
- `~/doknow-kb/audos/archive/*.md` — archived raw activity logs from the DoKnow build itself (several
  dated snapshots), more granular than the distilled blog posts. Worth at least one pass for anything
  the blog posts may have compressed away.

Read broadly first. Don't start synthesizing from the first few files.

## What to look for

Some starting hypotheses, drawn loosely from what's already in the evidence — treat these as things to
actively test against the material, not conclusions to confirm. You may find they're wrong, incomplete,
or that the real root causes cut differently:

1. **No real observability into a running job.** No way to cancel, no live log streaming, completion
   reports that truncate mid-sentence, no way to inspect intermediate state while a job runs — the
   opposite of a terminal-attached coding agent where every action streams to you as it happens.
2. **Job self-reports are structurally unreliable**, and nothing in the platform enforces or even offers
   automated verification before a job is allowed to claim success. This produced multiple false
   "verified" claims across both products.
3. **No transactional/rollback-safe schema operations.** A failed `CREATE TABLE` can leave a permanent,
   orphaned, unusable table with no cleanup tool — DDL that isn't atomic in an environment where an
   agent is expected to iterate quickly and imperfectly.
4. **Undocumented, surprising platform behavior discoverable only by trial and error** — silent
   `serial` ids where `uuid` was requested, injected columns, an integrations `proxy()` gateway that
   only works if you guess the right provider name and the right JSON-string body shape, with a silent
   200-OK-wrong-content failure mode if you get it wrong.
5. **Opaque publish/draft semantics.** Whether a build actually produced anything isn't reliably
   checkable — the file-tree API only reflects the *published* bundle, so "nothing was built" and "it
   was built but sits in draft" look identical unless you happen to know to check the Preview panel's
   Draft/Live toggle specifically.
6. **One-shot, non-recoverable credentials UX** — direct database access credentials can be generated
   exactly once, with no way to view or rotate them afterward if the string isn't captured at that
   moment, despite the backend's own error message implying a "regenerate" action exists.
7. **Unreliable scheduled/async execution** with no diagnostic tooling — recurring schedules
   inconsistently fire, root cause never isolated, no built-in way to check "did my scheduled job run,
   and if not, why."
8. **A chat-shell UI layer imposed by default**, that has to be actively "escaped" through a manual
   checklist (config flags, `useState` initializer ordering, chat-affordance behavior) rather than the
   platform simply producing a full-screen app when that's what was asked for.
9. **Workspace-shared mutable files** (a shared `Desktop.tsx`) mean a fix or bug in one app can silently
   affect every other app in the same workspace, which confounded at least one of our own experiments
   before we caught it — unclear app-level isolation boundaries.

Look for others we may have missed. Pay particular attention to *patterns that recur across multiple,
otherwise-unrelated findings* — that recurrence is the signal that you're looking at a root architectural
cause rather than an isolated bug.

## What "good" looks like for this audit

- Be honest and direct, not diplomatically hedged. If the evidence supports "this platform is not
  currently architected for hardcore agentic app development, full stop," say that plainly. If it
  supports a more qualified verdict, say that instead — but earn the qualification with evidence, don't
  default to one out of politeness.
- Ground every claim in specific, cited evidence — file paths, dates, job IDs, exact quotes where useful
  (this repo's own findings are written that way; match that standard).
- Explicitly call out where our own evidence is thin, contradictory, or where we corrected ourselves —
  the material contains several honest self-corrections (a "verified" claim later found to have missed a
  case, a filed bug later retracted, a retraction later found to itself be under-verified). Don't paper
  over those; they're part of an honest picture, not embarrassments to hide.
- Compare explicitly against how a terminal-attached coding agent (Claude Code, Codex) operates when it
  builds software: direct filesystem access under version control, a real shell it can run arbitrary
  commands and tests in, full unredacted output streamed live, the ability to inspect and modify runtime
  state directly, no black-box job dispatch, no ambiguity about whether a change was actually saved. Use
  that as the reference point for what "properly architected for this" would even mean — not as a sales
  pitch for either product, just as the working definition of the bar being measured against.

## Deliverable

A single audit report, structured roughly as:

1. **Executive verdict** — a few direct paragraphs answering the core question. Don't bury this at the
   end.
2. **Methodology** — what was reviewed, how much evidence exists, any real limits on what could be
   concluded from it.
3. **Top 10 architectural gaps, ranked by leverage** (the 80/20 — which fixes would unlock the most real
   capability). Clearly flag the top 3 as the most critical. For each: the pattern of evidence across
   multiple findings (cite them), the root cause as you understand it, what "fixed" would concretely look
   like, and — where relevant — how a terminal-attached coding agent avoids this class of problem
   entirely by construction.
4. **What's actually working well** — don't skip this. Several capability probes came back genuinely
   positive (native embedding generation once you find the right call shape, brute-force vector search
   performance, the base-URL shell-escape mechanism, library compatibility). An audit that only lists
   failures isn't credible or useful.
5. **Closing verdict**, restated plainly, with the single highest-leverage recommendation if Audos could
   only fix one thing.

Write it first as clean Markdown. Then produce a second, self-contained HTML version of the same
report — no external dependencies (inline any CSS), readable in both light and dark viewing contexts,
professional and readable rather than flashy. The tone throughout should match the rest of this repo:
direct, evidence-first, willing to say something doesn't work, equally willing to say something does.

Save both files under `docs/platform/reports/` in this repo, named
`architecture-audit-<today's date>.md` and the matching `.html`.
