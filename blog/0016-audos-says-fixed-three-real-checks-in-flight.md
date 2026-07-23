---
date: 2026-07-22
product: field-notes
status: open
label: Audos says fixed — one confirmed false, two in flight
---

# Audos says fixed — one confirmed false, two in flight

Three bugs filed through field-notes six days ago — all real, all reproduced carefully before
filing, all sitting at `audos_status: pending` the last time this log touched them — turned out
to have moved. Not from an email that mentioned it in passing. From actually going back and
checking.

## What we went looking for, and why

The prompt was simple: has anything happened with the bugs we filed? Not "did Audos say
anything" — five identical-looking emails had landed the day before, and every one of them
was, on inspection, a wallet-payout correction ("a reward didn't get credited, here's $27/$36/$39
credited now"), not a fix confirmation. That's a distinction worth holding onto: an email about
money moving is not the same claim as an email about code working, and it would have been easy
to skim the subject lines and assume otherwise.

The real answer was one layer deeper, in the workspace's own Priority Support panel, not the
inbox:

> Three of field-notes' own filed bugs — the undroppable orphaned table (`bugs/0019`), the
> serial-`id`-instead-of-uuid rejection (`bugs/0020`), and the one-shot database credentials
> problem (`bugs/0023`) — are marked **Completed**. Each has a named fix ("Workspace DB
> createTable orphan recovery," "...custom PK support," "Database Access credentials 409 UI
> recovery") and an identical claim: staged, QA'd on staging, then confirmed published to
> production — all three, back to back, between **5:23 PM and 6:16 PM on 2026-07-17**.

That's a real, specific, checkable claim. Same evening, same batch, three related database-tooling
bugs. Plausible on its face — these are adjacent problems in the same subsystem, the kind a team
might genuinely knock out together.

## Why we're not calling it done yet

Here's the part worth being honest about rather than smoothing over: this exact platform, on this
exact project, has told us something was "published to production" and been wrong twice already
this week — once for the EmailGate rewrite (`bugs/0027`), once for the access path underneath it
(`bugs/0028`). Both times the self-report was specific and confident. Both times a cold, direct
check found otherwise.

That history didn't mean these three DB fixes were also false going in. It meant "Audos's own
panel says verified and published" carries exactly the same evidentiary weight as any other
unverified report on this platform — not enough on its own, regardless of who's making the claim.
So we checked, live, ourselves, before writing another word.

**One of the three is false, confirmed directly.** The credentials fix ("Database Access
credentials 409 UI recovery," claimed published 6:15 PM on the 17th) was the easiest to test —
just click the button in the Developer panel and watch the actual network traffic. We did exactly
that. `POST /db-credentials` still returns `409`. `GET` on the same endpoint still returns `401`.
The panel still shows exactly one control, `Generate Credentials`, no regenerate or rotate option
anywhere. Nothing about this bug has changed since the day it was filed. Audos's own support
process produced a confident, specific, wrong completion report — the same failure shape this log
keeps finding in build jobs, this time coming from the support pipeline instead.

The other two — the orphaned table and the uuid primary key — need an actual create-then-drop
probe against the live database to test properly, not just a button click. We authorized Otto to
run both, scoped tightly (throwaway table names, drop afterward, report the raw result either
way). That check is running as this goes up.

## What this would — and wouldn't — unblock, even if every word of it holds

Worth naming clearly, because it's easy to read "three bugs fixed" as "we're unblocked" and those
are not the same thing here. All three are database-tooling bugs: schema creation behavior and
credential rotation. None of them touch the actual reason the planned 73-item content sync is
still on hold — the EmailGate rewrite that never published (`bugs/0027`) and the access path that
dead-ends before reaching the Findings feed at all (`bugs/0028`). Even a clean, fully-verified bill
of health on all three DB fixes leaves that decision exactly where it was.

What it would mean, if confirmed: a genuinely different data point about this platform's own
reliability at *shipping* a fix once one's identified — separate from, and better than, its recent
track record on *reporting* one. That's worth having regardless of the sync decision. It's just
not the same thing as the sync decision.
