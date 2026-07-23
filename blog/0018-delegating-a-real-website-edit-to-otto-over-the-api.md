---
date: 2026-07-08
product: Throughline
status: pass
label: Delegating a real website edit to Otto over the API, not the chat UI
---

# Delegating a real website edit to Otto over the API, not the chat UI

Everything this SDK log has documented about Otto so far has gone through the chat window — a human
typing into the workspace UI, watching a response come back. Throughline ran a different version of the
same question: does Otto's body-auth onboarding API (`/start` → `/status` → `/chat`, the same one from
`bugs/0029`'s auth saga) actually work as a real delegation surface — hand it a scoped, structured
brief, get a real edit made to a live marketing site, without a human driving the chat turn by turn?

## The brief, and what came back

The ask itself was specific on purpose: reposition the marketing site around podcast-first workflow and
VoicePrint as the differentiating engine, require human approval before anything goes live, use a
placeholder for the app CTA until a real URL was confirmed, and keep every forward-looking claim —
autonomous publishing, scheduling, analytics, social engagement, sales automation — explicitly fenced as
roadmap language, not shipped capability. Sent over the working API path, not typed into a chat box.

Otto routed it through `delegate_landing_page_edit` as a **draft, not an auto-publish** — task `78602`,
tracked and pollable, exactly the dispatch-then-poll shape this project's own SOP (`docs/platform/23`)
independently arrived at from the DoKnow side. A status check confirmed the task was running, unpublished,
still waiting on human review and a confirmed product URL before anything could go live.

## The part worth being honest about: the first draft wasn't good

Checking the actual copy that came back, the substance landed close to right — podcast-first framing,
VoicePrint front and center, roadmap language properly fenced. The visual execution didn't. A senior
design review called it plainly: heavy purple/blue gradient, pill navigation, emoji/card patterns,
over-large repeated blocks — the generic AI-SaaS look, not a considered brand.

> That's not a knock on delegation as an idea. It's what happens when "build me a landing page" is
> handed to an agent with no stronger creative direction than the brief itself — the same gravitational
> pull toward templated, AI-generated-looking design that shows up anywhere a model isn't given a
> specific point of view to work from. Worth remembering as a general lesson, not an Audos-specific one.

The fix wasn't to redo the delegation from scratch — it was to send a second, more specific brief:
reframe the direction as a premium editorial/podcast production desk, operational and human and
audio-native rather than pitch-deck-shaped, fix the CTA targets, add a proof-oriented section showing an
actual sample episode transformation. Sent through the same API. Otto came back having correctly
**superseded** the first draft — task `78611` explicitly replacing `78602`, carrying the original brief
plus every revision priority, still running, still unpublished, still gated on the same pre-publish
blockers (an approved lead-capture URL, no DNS changes made without asking).

## What this actually answers

Yes: a real content-production task can be delegated to Otto over the API and get a real, trackable,
non-destructive result — draft-gated, supersede-aware, safe by default rather than needing a human to
remember to check "did this auto-publish." What it doesn't answer, and shouldn't be read as answering,
is whether the *design judgment* on the other end of that delegation is any good without a human giving
it a specific enough point of view to work from. Both things were true in the same afternoon: the
plumbing worked cleanly, and the first draft still needed a real creative director's read before it was
worth shipping.

Source: `throughline/docs/decision-journals/2026-07-08-auto-api-website-execution-workstream.md`.
