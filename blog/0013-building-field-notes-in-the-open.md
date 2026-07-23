---
date: 2026-07-16
product: field-notes
status: pass
label: Building field-notes in the open
---

# Building field-notes in the open

Every entry in this log so far has been about DoKnow. Today a second thing got born: a small, real tool
called field-notes, meant to do exactly one job — publish this research, the bugs and feature requests
and reference docs, to a short list of people we personally approve, without depending on Notion,
Confluence, or any platform's own sharing model. Access control we control ourselves; content with real
structure instead of one more undifferentiated page in a folder.

We could have built it anywhere. We picked Audos on purpose, for the same reason DoKnow is here: since
we're already testing this platform's real capabilities, why not let building our own internal tool be
another one of those tests? Specifically, the sharpest version of a question this log has circled before
without ever fully answering: can a brand-new Audos workspace be born already outside the chat shell —
full-screen from the very first paint, no retrofit required — or does that only ever work because we'd
already fixed a shared file in some other workspace first?

That question mattered enough to get its own honest correction before the day was half over. The first
instinct was to build field-notes as just another app inside DoKnow's own workspace — fast, and every
capability we'd already proven (the OpenAI proxy, the JSON-array search) would already be sitting there,
no re-verification needed. Wrong call, and said so plainly: DoKnow's workspace and field-notes' workspace
needed to be genuinely separate things, not because Audos required it, but because conflating a product
workspace with internal tooling infrastructure was the kind of shortcut that quietly makes a mess later.
A real, fresh workspace got created. Its own `Desktop.tsx` has never been touched by any prior fix — which
is exactly what made it the first legitimate chance to run the experiment for real.

We designed the data model before writing anything against the platform, on purpose — one table,
`content_items`, with a `content_type` column doing the work five separate tables would otherwise need,
because the existing build-log script already proved blog posts, bugs, and feature requests are one
shape with different metadata, not five different shapes. Writing the schema down first was itself part
of the test: does giving an agent-driven platform an explicit, complete contract up front produce a more
reliable result than letting it improvise structure as it goes?

The first real answer came back mixed. Creating the tables surfaced a genuine platform bug: a failed
foreign-key constraint didn't roll back, it left a physical, orphaned table sitting in the database that
blocked every future attempt to use that name, with no tool available to clear it. We also learned
something the schema had gotten wrong on our own end — Audos always generates a table's `id` as its own
auto-increment integer, never the `uuid` we'd assumed, which meant every foreign key in the design needed
correcting before anything else could proceed.

The second attempt, a combined job to fix the table problem and actually build the app, ran for close to
two hours. Along the way it found something else worth flagging carefully: a possible gap in the raw-SQL
restriction, where the validator might only check the first keyword of a query and not the whole
statement. We're not calling that confirmed — the job tried it, gave up, and worked around the orphaned
table a different way instead, which is itself evidence the bypass probably didn't pan out. Flagged for
Audos to check on their own, not asserted as a working exploit based on one incomplete attempt.

Then, per its own account, the job actually succeeded at the thing we set out to test. It built the app,
edited the workspace's landing configuration, patched the shared shell file, and reported testing all of
it live: the bare workspace URL loading the password gate full-screen with no chat and no dock, a
complete sign-in-to-feed-to-sign-out flow working end to end, a reload landing a returning visitor
straight back on the feed with no flash of old UI in between. Read on its own, that's the answer to the
question this log has been circling: yes, briefed upfront, from a genuinely fresh workspace, it can be
born outside the shell from the very first paint.

Except the job crashed on an unrelated cosmetic fix before it ever posted that as a finished report — and
when we went and checked the actual file tree ourselves afterward, it was completely empty. No app file.
No configuration change. Nothing. Whatever got tested, live, inside that job's own run, didn't survive
the run not finishing cleanly.

That's not the same finding as "the checklist doesn't work." It's a different, in some ways more useful
one: the platform's database writes are durable the moment they happen, ordinary transactions, already
confirmed and still sitting there when we checked independently. File writes are not — they appear to sit
uncommitted until some later step a job may never reach, and an error before that point silently discards
everything, even work the job itself had just finished watching succeed. We now have a standing rule out
of this, the same discipline this whole log keeps re-learning in different forms: never trust a job's own
account of what it verified as durable until you've gone and checked yourself, independently, after it's
actually done.

So, honestly, where this stands tonight: the database side of field-notes is real, confirmed by hand, not
by a job's word for it — the content table, the access table, a sessions table (recreated under a new
name, since the original one is now permanently stuck behind the orphan bug), one working test
credential. The actual application doesn't exist yet. The shell-escape question is neither answered yes
nor answered no — the one attempt that would have answered it was destroyed by something else entirely
before the result could be checked.

Four new findings came out of today, on top of everything already tracked: the orphaned-table bug, the
undocumented serial-id behavior, the unconfirmed raw-SQL gap, and the file-durability issue, all filed in
the same bug tracker as the rest of this build. Then we filed all four with Audos directly, and one of them
— the file-durability finding — got pulled to a human engineer within minutes of submitting it.

> **Update, same evening.** The retry ran. It came back clean, and something in its report didn't add up
> in an interesting way: it said the previous job's work — the app, the config, the hooks — had already
> been there when it started. Not rebuilt. Found.

That's the kind of claim you don't just write down. We went and checked ourselves — not the API endpoint
we'd checked before, but the workspace's own Preview panel, Draft against Live, side by side. Live showed
the same tired old marketing page it always had. Draft showed a real password gate, styled, waiting. We
typed in the test credentials ourselves and watched it open onto an actual content feed — three sections,
each with a plain, correct empty state. We clicked the assistant button ourselves and watched a small
popup open over the still-visible feed, not the old chat interface swallowing the screen.

The app was never gone. It had been sitting in an unpublished draft the entire time, and the tool we used
to check for it — `GET /api/space/{id}/files` — only ever shows what's been *published*, never what's
sitting in draft. We checked a doorway that only opens for finished work and concluded the room was
empty. It wasn't empty. We just hadn't opened the right door.

That means the file-durability bug we filed a few hours earlier wasn't a bug. We went back to the same
support thread — the one with a human already reading it — and said so plainly: don't spend time on this,
we were wrong, here's exactly how we were wrong, here's what to fix instead (the fact that a plain file
check can't tell "nothing built" apart from "built but only in draft," which is a real if smaller gap).
That felt like the more important thing to get right today, more than the original finding itself: when
you're wrong on the record, in front of someone whose time you asked for, you say so in the same place,
not quietly somewhere else.

So, the actual answer to the question this whole entry opened with: yes. Briefed upfront, in a workspace
that had never been touched by any prior fix, field-notes was born outside the chat shell from its very
first paint. Not claimed — checked, by hand, twice, the same evening we almost got it wrong.

> **Update, later the same evening.** Turned out even that retraction wasn't as clean as it felt at the
> time.

Getting the retraction right meant checking the app's *current* state ourselves — logging in, reaching
the real feed, clicking the popup. What it didn't mean was independently checking the *intermediate*
state: whether the first job's work had actually survived, untouched, in the gap before the second job
started. That claim — "the previous job's work was already present, not rebuilt" — came from the second
job's own self-report, which is exactly the kind of thing this whole log has spent two weeks learning not
to trust. We'd applied the discipline once and then quietly skipped applying it to our own correction. By
the time the gap was noticed, the second job had already run — whatever the true intermediate state was,
it's gone now, unrecoverable. We're not filing a further correction over a claim we can no longer check;
the existing retraction stands, with this gap written down plainly instead of pretended away. A new
standing rule came out of it: nothing gets posted to Audos support anymore unless the underlying claim is
confirmed with full independent certainty first — verifying our own verification, not just the platform.

Then came the part that actually mattered for finishing field-notes: with the app confirmed working, the
last piece was getting real content into `content_items`. The direct route — Developer panel, Database
Access, "Generate Credentials," connect with a normal Postgres client, skip the AI-build-job path
entirely since inserting rows isn't an app build — hit a wall immediately. Credentials had already been
generated once earlier in the session and the connection string was never captured. Every click since
returned a flat `409 Conflict`: *"Credentials already exist. Use regenerate to rotate them."* A full scan
of the panel's interactive elements found exactly one control tied to the feature — `Generate
Credentials` — and nothing else. No regenerate, no rotate, no view. The error message references a
capability the interface simply doesn't expose.

Before filing that as a bug, we asked Otto directly, in the workspace's own chat: any way to view or
reset these from in here? Otto checked its own API surface and came back straight: no, on both counts.
Its database tools run through a separate, short-lived token unrelated to this specific credential, and
— more usefully — it confirmed there's genuinely no "view" path for anyone, since the string is only ever
shown once at generation and isn't kept in plaintext anywhere retrievable. That's an independent
confirmation from the platform's own side, not just our own reproduction of a stuck button, which is what
actually cleared the bar for filing it. Otto prepared the report and submitted it as Priority; the
confirmation came back with an unexpected detail attached — a $27 bug bounty, payable to the workspace
wallet if the fix turns out to be real. We went looking for where that gets tracked and found nothing:
Wallet showed a $0.00 balance, no pending line item, "Investment Decisions: Not Yet Funded." The bounty
exists only as one-time text in a confirmation dialog, nowhere else. Small, but worth naming as its own
gap — a commitment with no visible state to check later.

`content_items` is still empty tonight. Two paths remain open for actually getting content in — a
narrowly-scoped build job that only inserts rows, or waiting on the credentials fix to land — and neither
has been chosen yet. That's fine. The point of building field-notes as its own workspace was never to
finish it fast; it was to have a place to run the platform's sharper edges through their paces without
touching DoKnow while doing it. Today it did exactly that: one real, if narrow, capability bug found and
filed with independent confirmation behind it; one big open question about the platform finally answered
by hand, twice; one lesson about verifying our own corrections as carefully as we verify the platform.
DoKnow itself never had to wait on any of it.
