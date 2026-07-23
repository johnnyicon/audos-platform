---
date: 2026-04-05
area: ai-api
status: confirmed
label: The AI hook is hard-locked to gpt-4o-mini-2024-07-18 with no model selection, and streaming isn't supported
---

**Hypothesis:** Given the AI hook is described in Audos's own SDK docs as a general-purpose generation
primitive, there should be some path — a parameter, a workspace setting — to select a more capable model
for tasks that actually need one (complex generation, voice-fingerprint training), and some way to
stream a response rather than wait for the whole thing.

**Method:** Read Audos's own SDK documentation package (SDK-00 through SDK-08) end to end and asked the
two direct follow-up questions this raised.

**Result:** Confirmed, plainly, from Audos's own docs: the model is locked to
`gpt-4o-mini-2024-07-18` (128k context), with **no model-selection path exposed** — no parameter, no
workspace setting, nothing. Streaming is **not supported**; the SDK's own guidance is to use loading
states instead of expecting incremental output. Both are real ceilings for tasks like voice-fingerprint
training or complex multi-step generation, not just a preference — there's no documented or undocumented
path around either one, unlike some of this project's other findings where a workaround existed once you
knew to look.

Source: `throughline-forge/tmp/handoff-nicholas-audos-sdk.md`.
