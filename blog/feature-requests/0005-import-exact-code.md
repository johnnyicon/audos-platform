---
date: 2026-07-13
priority: 3
status: not filed
label: A genuine, mechanical "import this exact code" path
---

Getting our own hand-built UI to render faithfully required pasting a large HTML/CSS/JS file into a chat
message with the instruction "port this verbatim, do not redesign, restyle, or improve anything" — stated
explicitly, because compliance is an LLM judgment call, not a mechanical guarantee. It worked every time
we tried it, but it worked because we got lucky with agent compliance, not because there's a real feature
for it.

A literal file-upload → deterministic-wrap → publish path — no LLM interpretation of layout or content,
only the boilerplate React wrapping that's genuinely mechanical — would make this a guarantee instead of
a bet, and would remove the one real near-miss we hit doing exactly this migration: a build agent's first
pass claimed a verbatim port but had only pasted the CSS literally, reconstructing the HTML/script from a
description instead.
