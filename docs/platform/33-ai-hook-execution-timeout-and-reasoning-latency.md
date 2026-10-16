# The AI hook's ~40 s execution cutoff, and why reasoning models hit it

*Written 2026-10-15 from the Throughline script-loop build (`throughline`, workspace-351699, hook
`ai-api`). Live evidence: `throughline/docs/plans/script-loop/evidence/replay/2026-10-15-live/`.*

## The limit

The AI hook (`POST /api/hooks/execute/workspace-{id}/ai-api`, action `chat`) is a server function,
and a server function has the runtime's execution limit (`06-capabilities-reference.md` says ~30 s;
`08-server-function-runtime.md` says the same). **Observed on 15 Oct 2026: the hook gives up at
about 40–45 s** and returns HTTP 500 with this body (durations measured by the hook itself):

```json
{"error":"No response from AI model","_meta":{"success":true,"durationMs":41481,"logs":[]}}
{"error":"No response from AI model","_meta":{"success":true,"durationMs":44551,"logs":[]}}
```

`_meta.success` is `true` and `logs` is empty: the hook ran fine, the model call under it was cut
off. There is no retry, no partial text, and nothing in the body says it was a timeout. The
upstream request may well have completed on OpenAI's side after the hook stopped waiting; you pay
for it either way.

This is a per-request wall-clock limit on the hook. It is **not** a token limit (Quirk 6 in
`docs/AI-HOOK-CAPABILITY-MATRIX.md`) and not an OpenAI timeout.

## Why reasoning models hit it

`gpt-5.5` (and GPT-5 / o-series generally) reasons before it writes, then writes. A request that
asks for a long structured answer, such as a JSON list of 25 suggestions over a 400-word script
section with 26 house rules to apply, took **20–55 s per request** in our measurements, with the
longest ones cut off. The same prompt producing a short answer returns in 2–8 s. So the time is
spent on **output length plus reasoning**, not on input size; a request that reads a lot but
answers little is fast.

Measured per-request times, `gpt-5.5`, `maxTokens: 4000`, `temperature: 1`, one script section
per request (Throughline EP5 draft, 1,072 words, 9 sections):

| Prompt kind | Per request | Notes |
|---|---|---|
| Facts (sources attached, few findings) | 8–13 s | never cut off |
| Wording (plain-language rewrites) | 10–40 s | the 40 s one was close |
| House rules (26 rules applied, many findings) | 6–55 s | 2 of 9 cut off and retried |

## What to do about it

1. **Keep every request well under 30 s of generation.** Split the work by section or by task, not
   by one big prompt, and cap the answer length you ask for (we moved from "at most 25
   suggestions" to "at most 10" per request).
2. **Run the pieces in parallel.** The hook handles concurrent requests: three requests sent at once
   each answered in 2 s, wall-clock 2 s total. A nine-section check that took 303 s sequentially took
   124 s with four requests in flight at a time, and would be bounded by the slowest single request
   with more.
3. **Retry a cut-off once**, with the same input. The retry usually succeeds; it is the variance of
   reasoning time, not a deterministic failure.
4. **Do not use `maxTokens` to speed it up.** Below ~4000 a reasoning model returns empty text with
   `finish_reason: length` (Quirk 6). The length you control is the length you *ask* for in the prompt.
5. **For long reads with short answers, the limit is fine.** For long answers, consider a non-reasoning
   model (`gpt-4.1`) and compare quality, or a direct provider call outside the hook.

## Two related corrections, same day

- **`temperature` must be sent as `1` for `gpt-5.5`.** The hook's default (0.7) is forwarded and
  refused upstream: `400 "Unsupported value: 'temperature' does not support 0.7 with this model. Only
  the default (1) value is supported."`, surfaced as the hook's HTTP 500 with the OpenAI error in
  `_meta.logs[0]`. The matrix's parameter table (`temperature`, default 0.7) does not say this.
- **The stock `ai-api` hook on workspace-351699 serves `gpt-5.5` today** (served as
  `gpt-5.5-2026-04-23`, with `reasoning_tokens` in `usage`), so Quirk 4's "requires the corrected
  `ai2-api` hook" is no longer true for this workspace as of 15 Oct 2026. Whether Audos updated the
  stock hook everywhere is not verified; check `_meta.logs[0]` on a 500 before assuming.

## Open with Audos

Backlog `#45`: the hook's execution limit is undocumented for AI calls and too short for reasoning
models' long answers; ask for a longer limit on the AI hook (or a documented one), and for a
`timeout` field in the error body instead of "No response from AI model". Draft to file:
`bug-reports/2026-10-15-ai-hook-40s-cutoff.md`.

## Rate limiting of many small requests (found 15 Oct 2026, backlog #46)

The cutoff above bites long requests. Many **short** requests hit a different, undocumented limit:

- A burst of ~24 requests in 11 s makes the hook refuse **every** request for about a minute: HTTP 500,
  body `{"error":"AI generation failed"}`, returned instantly. Not a timeout, not an upstream error.
- At a steady 20 requests a minute it still refuses 10–30% at random (24 of 241 in one run, 87 of 288
  in another, same hour); at 30 a minute, about one a minute. The same request succeeds on retry.
- Nothing distinguishes a refusal from a real model failure: no 429, no `Retry-After`, same body.

Measured with an agentic checker (≈55 requests a check, four conversations in parallel, 1–5 s each).
The one-shot design (11–29 requests of 10–40 s per check) saw none of this on the same hook the same day.

**What to do:** prefer fewer, longer requests. If you must loop, pace to at most 20 requests a minute
per workspace, retry an instant "AI generation failed" up to three times with growing pauses (3, 6, 12 s),
and log it as a refusal, not a model error. Evidence: Throughline
`docs/plans/script-loop/evidence/replay/2026-10-15-agentic/README.md`.
