---
date: 2026-07-23
product: DoKnow
status: pass
label: We asked Otto directly what Audos can do that our capability matrix hadn't captured — and then made it prove every claim
---

# Asking Otto what else Audos can do

The capability matrix this SDK maintains (`CAPABILITY-MATRIX.md`) is synthesized entirely from our own
hands-on findings — six months of bugs, experiments, and feature requests, one row per capability. It's
thorough for what we've actually built and broken. But it only covers what we happened to run into. So we
tried something more direct: hand Otto the current matrix, over the same onboarding-API chat path we've
used since the very first auth investigation (`bugs/0029`), and ask it point-blank what's missing.

## The answer was long. That was the problem.

Otto came back with a real list — Stripe payments, a secrets proxy for BYOK API keys, a full Meta ads
pipeline, video/voiceover/music generation, analytics and dashboards, CRM tooling, browser automation,
workspace snapshots, and more. Useful, but it's exactly the shape of information this project has learned
not to trust at face value. Otto describing its own capabilities is a self-report, the same category as a
build job claiming "Complete" while nothing shipped (`docs/platform/23`). A longer, more confident
description doesn't upgrade a claim from 📄 (Audos says so) to ✅ (we checked) — only actually running the
thing does.

So before adding a single row to the matrix, we split Otto's list by one question: **does testing this
have a side effect?** Querying analytics, researching keywords, listing campaigns, previewing ad copy —
none of that touches a real person or spends real money, so we ran all of it for real. Launching an actual
ad campaign, sending real Instagram DMs, spending on Stripe live mode — those we left alone. Video, image,
and voiceover generation sit in between: no one gets hurt, but they burn real workspace credits, so we
checked first before spending any.

## What we actually found

**Analytics held up, mostly.** `query_analytics` and `query_events` both came back with real, plausible
numbers, and the contact count cross-checked cleanly against a second, independent tool call — a genuine
corroboration, not the same figure echoed twice. But a more general query tool, `query_data_source`,
returned a degenerate, all-blank result for event data in the same workspace, in the same conversation,
where `query_events` had just reported 187 real events. Two tools that should agree, didn't.

**Ads tools were real and useful, with one sharp edge.** Keyword research returned genuine Google Ads
metrics with the kind of row-to-row variation that doesn't come from a template. Ad copy generation wrote
three variants that actually referenced DoKnow's real pain point — forgotten reading, not generic filler —
and did it cold, with no prior campaign history. But Meta's geo-targeting tool silently failed on
`"Austin, TX"`, the single most natural way to type a location, and only worked once we switched to a bare
city name or a ZIP code.

**Media generation was the most fun to verify, and the least trustworthy to leave un-checked.** We asked
Otto to generate a real image, a real 8-second video (Veo3), and a real voiceover clip (ElevenLabs) — then
downloaded every file ourselves and inspected it directly rather than trusting the returned URL. The image
rendered exactly the text we asked for. `ffprobe` confirmed the video was a real, playable 8-second MP4
and the voiceover a real 8.6-second MP3. All three were genuine. Along the way we also found a real bug:
asking for the full list of available voices returns about 85 of them, but 82 have a literal `undefined`
where their ID should be — only the three hardcoded fallback voices actually work.

We also tried to price it out. Wallet balance moved by exactly three cents across the whole session, and
Otto could trace that three cents to a single, cleanly-itemized line: ElevenLabs TTS, billed per character.
The image and video generations, by contrast, both drew from a shared "AI usage" bucket with no
per-generation number attached at all — real cost, real spend, just not one Audos's own tools can itemize
for you today.

## What changed in the SDK

Every capability above got its own file — `docs/platform/capabilities/{analytics,ads,media}/*.md` — plus
three short index docs, following the same progressive-disclosure shape as the `otto-pilot` skill's own
`references/` folder: cheap to skim the index, load only the one capability file you actually need. A new
architectural note (`docs/platform/29`) records the thing that made all of this necessary in the first
place: the tools Otto called for us in this whole pass are Otto-chat-orchestrated, not `platform.*` hooks
an app you build could call directly — a distinction worth having explicit, since the marketing language
for both sounds identical.

The capability matrix picked up two new categories (Analytics & Reporting, Ads & Marketing) and a Media
Generation section, with sources that say plainly which rows we ran ourselves and which are still Otto's
word. That's still the whole discipline this file exists to enforce — a claim doesn't turn green until
someone actually checks.

Source: Otto chat, DoKnow workspace (`8a65a4ac-5a22-435f-b55f-c41ea34ca00d`), 2026-07-23.
