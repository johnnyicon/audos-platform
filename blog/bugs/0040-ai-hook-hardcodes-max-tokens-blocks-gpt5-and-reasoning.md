---
date: 2026-08-10
area: ai-hook
status: open
filed: no
label: The ai-api hook hardcodes max_tokens, blocking every current-generation OpenAI model (GPT-5, o-series) that are otherwise authorized
---

The `ai-api` hook always sends OpenAI's `max_tokens` parameter. OpenAI's current-generation models —
GPT-5 (all sizes), `gpt-5.5`, `gpt-5.6`, and the o-series (`o1`, `o3`, `o3-mini`, `o4-mini`) — reject
`max_tokens` and require `max_completion_tokens`. So every one of these models fails through the hook,
even though the proxy account is authorized for them.

**How we know they're authorized (not just missing).** A 43-name live audit read the *upstream* error out
of `_meta.logs`. These models return:

```
400 {"error":{"message":"Unsupported parameter: 'max_tokens' is not supported with this model.
     Use 'max_completion_tokens' instead.","type":"invalid_request_error","code":null}}
```

That is a **parameter** error, which OpenAI only returns *after* it has resolved the model and authorized
the key. An unavailable model returns `404 model_not_found` instead — which is exactly what `o1-mini`,
`gpt-5.6-sol`, and every `claude-*` / `gemini-*` name return in the same sweep (the negative control). So
the 400-vs-404 split cleanly separates "reachable but blocked by the wrong parameter" from "genuinely not
available."

**Confirmed it's the hook, not the caller.** Every attempt to work around it from outside failed:
`maxTokens: 0` still sends `max_tokens: 0`; `maxTokens: null` produces a type error; passing our own
`max_completion_tokens` in the body is ignored (the hook adds `max_tokens` alongside and drops unknown
fields). There is no request shape that reaches these models through the current hook.

**Impact.** The hook can only reach the GPT-4/3.5 generation. Every modern OpenAI model — including the
only ones that expose reasoning/effort — is unreachable for app code, despite being paid-for and available
on the proxy. For any workspace building on the in-app AI, this silently caps them a full model generation
behind.

**Scope / provenance.** This is the behaviour of one workspace's `ai-api` server function (it's
workspace-built, not platform-provided — `docs/platform/29`). We have one hook sample and can't read its
source, so we can't confirm whether Audos's own generated hooks share the code. Stated as a hook-code
bug, not a platform defect.

**Fix.** In the hook, branch on the model id: send `max_completion_tokens` for `gpt-5*` and `o[1-9]*`,
`max_tokens` otherwise. A one-line conditional. The runtime call is
`platform.integrations.proxy('openai', '/v1/chat/completions', { … })` (per `CAPABILITY-MATRIX.md`'s
embeddings row), so a corrected or self-built hook has everything it needs.

## Update 2026-08-14 — fix verified working

The one-line fix is confirmed. A corrected hook (`ai2-api`) was created in the DoKnow workspace via a
Cursor delegation job (#108698), branching the token parameter on model id. Verified by independent calls
rather than the job's self-report: `gpt-5` → `gpt-5-2025-08-07`, `gpt-5-mini`, `o3` → `o3-2025-04-16`, and
`o4-mini` all return real generated text, on **Audos's own OpenAI key** via
`platform.integrations.proxy('openai', …)`. Source: `sdk/hooks/ai-api.reference.js`.

So this is no longer "a confirmed block with a likely fix" — the fix is proven. The stock `ai-api` hook
remains broken, and its source can't be read or written by any tool available to a builder, which is why
the remedy had to be a second hook rather than an edit. **That's the part still worth Audos's attention:**
if `ai-api` is scaffolded for every workspace, every builder on the in-app AI is silently capped a
generation behind, and cannot fix it themselves.

**One caveat discovered while verifying** (now its own finding): reasoning models spend the token budget
on *internal* reasoning before emitting visible text. `gpt-5` at `maxTokens: 2000` returned empty text with
`finish_reason: length` — all 2000 tokens went to reasoning. At 4000 it produced real output. Anyone
implementing this fix should raise the default budget for reasoning models, or callers will see silent
empty responses.

Source: live audit against `workspace-351699/ai-api`, 2026-08-10 (see doknow-kb
`research/audos-ai-capabilities.md`).
