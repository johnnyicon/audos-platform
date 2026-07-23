---
date: 2026-07-12
product: DoKnow
status: shipped
label: Kickoff
---

# Kickoff: standing up DoKnow on Audos

DoKnow is our test case for how far Audos can actually go. The product itself: turn the videos,
articles, podcasts, and PDFs you pile up — or just a topic you name — into a real, leveled course,
then coach you through it in 5–10 minute daily lessons until it sticks. The full spec lives in our
`doknow-kb` repo; the differentiator is a proactive Coach, not the content generator.

We onboarded through Audos's agent API — email, a 4-digit OTP, and then an autonomous 11-stage build:
customer research, problem mapping, an AI tool suite, brand identity, a hero video, the workspace
"space," the landing page, and ad creatives. Total time: about six minutes.

What came back was genuinely good on the surface. Audos picked a teal-and-coral palette, wrote its own
tagline ("From what you save to what you know"), and correctly identified the target customer —
curious, busy people who save more than they finish — almost word for word from our own product plan.
Notably, of everything it could have built first, it chose to scaffold a "Coach Queue" app — the exact
piece we'd called the product's differentiator. That's a good sign about how well the onboarding
research actually reads a business idea.

The honest gaps showed immediately too. The hero copy was still generic template language ("Become
who you set out to be") — only the tagline was ours. And the workspace shipped with exactly one app,
not the suite the plan called for. We rewrote the hero, cleaned up the CTAs, and then built two more
apps — a Course Builder and a Lesson Player — to round out the core loop.

First real platform lesson: **build jobs run on Cursor Background Agents**, and that's Audos's own
shared account, not the workspace owner's. When it's over its usage limit, every job instant-fails
with `usage_limit_exceeded` — and there's no way to tell from inside the workspace whose account that
is or how long it'll take to recover. It cleared on its own within the session, but it's a real
dependency to know about if you're planning a build.

By the end of the first pass we had a live landing page, three working apps, and a workspace that
looked like a real product. It also raised the question that shaped everything after: does it *feel*
like a product, or does it feel like a chatbot with some tools bolted on? That's where the next post
picks up.
