# Audos AI Hook — Capability Matrix

**Endpoint:** `https://audos.com/api/hooks/execute/workspace-351699/ai-api`  
**Method:** `POST`  
**Content-Type:** `application/json`  
**Authentication:** None (workspace URL is public)

---

## Critical Platform Quirks

These two behaviors are non-obvious, cost time to discover, and affect every caller.

> **Major correction, 2026-08-10.** A 43-model live audit (reading the *upstream* error out of
> `_meta.logs`) overturned two earlier claims in this doc: **vision works**, and **GPT-5 + the o-series
> are reachable** — they were mis-recorded as unavailable. See "Available Models" and the two new quirks
> below. The old "accepted but broken" framing conflated three genuinely different cases.

### Quirk 1: OpenAI proxy only — not a multi-provider gateway

The hook is an OpenAI proxy. It accepts any model string without pre-validation, but only OpenAI-served
IDs work. There are **three** distinct outcomes, not two — and the difference is only visible in
`_meta.logs` (see Quirk 3):

1. **Works** — OpenAI GPT-4/3.5 family. Returns text.
2. **Reachable but hook-blocked** — GPT-5 (all sizes), `gpt-5.5`, `gpt-5.6`, `o1`, `o3`, `o3-mini`,
   `o4-mini`. Upstream returns `400 "Unsupported parameter: 'max_tokens' … Use 'max_completion_tokens'"`.
   That is a *parameter* error, which OpenAI only returns for a model it has resolved and authorized — so
   these models **are available**; the hook just sends the wrong token parameter (see Quirk 4).
3. **Genuinely unavailable** — `404 model_not_found`. Non-OpenAI models (Claude, Fable, Gemini, DeepSeek,
   Kimi) and models this account lacks (`o1-mini`, `gpt-5.6-sol`). No hook change reaches these through an
   OpenAI proxy.

The 400-vs-404 distinction is the discriminator: `400` param error = reachable, `404` = not.

### Quirk 2: `generate` action has a hardcoded 1000-token output cap

The `generate` action ignores whatever `maxTokens` you send. Output is capped at exactly 1000 completion tokens regardless of the value in the request. Verified across four probes:

| `maxTokens` sent | `completionTokens` returned |
|------------------|----------------------------|
| 200 | 1000 |
| 500 | 1000 |
| 2000 | 1000 |
| 8000 | 1000 |

The `chat` action does not have this cap. The same prompt via `chat` with `maxTokens: 8000` returned 7816 characters (~1900+ tokens).

**For any structured output that may exceed 1000 tokens, use `chat` instead of `generate`.** This includes arc generation, show notes, research synthesis, and any multi-section structured content.

### Quirk 3: the real error hides in `_meta.logs` (2026-08-10)

The top-level `error` field now flattens every model-call failure to a generic `"AI generation failed"`
(HTTP 500). The **real upstream reason is intact** — it moved into `_meta.logs`:

```json
{ "error": "AI generation failed", "status": 404,
  "_meta": { "logs": ["[ERROR] OpenAI proxy error: 404 {\"error\":{\"message\":\"The model `X` does not exist…\"}}"] } }
```

Always read `_meta.logs[0]` to see why a call failed. The `status` field inside the body is the *upstream*
OpenAI status (404 = bad model, 400 = bad parameter), distinct from the hook's HTTP 500 wrapper. The
`@audos/sdk` AI client surfaces this line automatically as the thrown error's `serverError`.

### Quirk 4: `max_tokens` blocks current-generation models — **FIXED AND VERIFIED (2026-08-14)**

GPT-5 and the o-series require `max_completion_tokens`; the original `ai-api` hook always sends
`max_tokens` (even `maxTokens: 0` sends `max_tokens: 0`; passing your own `max_completion_tokens` is
ignored). Filed as `blog/bugs/0040`.

**This is now solved.** A corrected hook, `ai2-api`, was created in the DoKnow workspace
(`workspace-156396`) via a Cursor delegation job. It branches the token parameter on model id. Verified
live by independent calls, not the job's self-report:

| Model | Served as | Result |
|---|---|---|
| `gpt-5` | `gpt-5-2025-08-07` | ✅ real output |
| `gpt-5-mini` | `gpt-5-mini-2025-08-07` | ✅ |
| `o3` | `o3-2025-04-16` | ✅ |
| `o4-mini` | `o4-mini-2025-04-16` | ✅ |

Endpoint: `POST https://audos.com/api/hooks/execute/workspace-{id}/ai2-api` (the hook is per-workspace and
must be created there first — see `sdk/HOW-TO-USE-AUDOS-AI.md` §9). It calls
`platform.integrations.proxy('openai', …)`, so it runs on **Audos's OpenAI key** — no BYOK needed.
Source: `sdk/hooks/ai-api.reference.js`.

**The original `ai-api` hook is still broken** and its source is unreadable/unwritable, which is why the
fix is a new hook alongside it rather than an edit.

### Quirk 6: reasoning models silently return empty text if `maxTokens` is too low (2026-08-14)

Reasoning models (GPT-5, o-series) spend part of the budget on **internal reasoning** before emitting any
visible text, and that spend counts against `max_completion_tokens`. Measured on `gpt-5` with a
moderately complex prompt:

| `maxTokens` | completion tokens | reasoning tokens | visible text |
|---|---|---|---|
| 2000 | 2000 | 2000 | **0 chars** (`finish_reason: length`) |
| 4000 | 1746 | 1152 | 2266 chars ✅ |
| 8000 | 1926 | 1408 | 1892 chars ✅ |

So a budget that works fine for GPT-4 can yield **nothing** on GPT-5. Use ≥4000 for reasoning models, more
for complex prompts. The `@audos/sdk` AI client handles this: it defaults reasoning models to 4000 and
raises an explicit "raise maxTokens" error instead of returning an empty string.

### Quirk 5: vision works (2026-08-10)

Image understanding is available on `gpt-4o`/`gpt-4.1`. Pass `messages[].content` as an OpenAI-style
array with an `image_url` part (public URL or `data:` URI). Verified: a red PNG returned "Red". The
earlier "vision not exposed" note was wrong (the test image had been malformed).

---

## Overview

The Audos AI hook is an OpenAI proxy with a normalized request/response layer. It accepts requests with model names, forwards them to OpenAI's API, and returns a normalized response shape.

**Key behaviors:**
1. Model name is not pre-validated — any string is accepted
2. Non-OpenAI model names route to OpenAI and 404 silently (via `generate`) or visibly (via `chat`)
3. The `generate` action caps output at 1000 tokens regardless of `maxTokens`
4. The `chat` action honors `maxTokens` as expected

---

## Available Models

Three tiers, from the 43-name audit on 2026-08-10 (status by *upstream* status code):

**Tier 1 — working now** (returns text):

| Model | Served as | Notes |
|-------|-----------|-------|
| `gpt-4o` | `gpt-4o-2024-08-06` | vision-capable |
| `gpt-4o-mini` | `gpt-4o-mini-2024-07-18` | **default**; vision-capable; lower cost |
| `gpt-4.1` / `-mini` / `-nano` | `gpt-4.1-*-2025-04-14` | vision-capable |
| `gpt-4-turbo` | `gpt-4-turbo-2024-04-09` | |
| `gpt-4` | `gpt-4-0613` | |
| `gpt-3.5-turbo` | `gpt-3.5-turbo-0125` | |

**Tier 2 — frontier + reasoning, WORKING via the corrected `ai2-api` hook (Quirk 4)** — authorized on the
proxy; blocked only by the *original* hook's wrong parameter:

`gpt-5`, `gpt-5-mini`, `gpt-5-nano`, `gpt-5.5`, `gpt-5.6`, `o1`, `o3`, `o3-mini`, `o4-mini`.
Against the old `ai-api` they return `400 "use max_completion_tokens"` (not `404`, which is what proved
they were available all along). Against `ai2-api` they generate normally — see Quirk 4 for the verified
results, and Quirk 6 for the token-budget gotcha.

**Tier 3 — genuinely unavailable** (`404 model_not_found`): all `claude-*`, `fable-5` /
`claude-fable-5*`, all `gemini-*`, `deepseek-chat`, `moonshotai/kimi-k2.5`, `o1-mini`, `gpt-5.6-sol`,
`chatgpt-4o-latest`, and `anthropic/…` / `google/…` prefixes. Claude/Fable run on the **Otto/job
surface**, which is not an app-callable API.

**For text today, use a Tier-1 model** (`gpt-4o`, `gpt-4o-mini`, `gpt-4.1`, `gpt-4-turbo`). For GPT-5 or
reasoning models, a corrected hook is required (Quirk 4).

---

## Request Shape

### Actions

The hook supports two actions:

#### `generate` — Single-turn text generation

```json
{
  "action": "generate",
  "prompt": "Write a 2-sentence caption for a podcast episode",
  "model": "gpt-4o",
  "systemPrompt": "You are a social media expert. Be concise.",
  "maxTokens": 100,
  "temperature": 0.7
}
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `action` | string | Yes | — | Must be `"generate"` |
| `prompt` | string | Yes | — | The user query or instruction |
| `model` | string | No | `gpt-4o-mini` | Model to use. Must be a valid OpenAI model ID. |
| `systemPrompt` | string | No | — | System instructions to guide the model |
| `maxTokens` | integer | No | 1200 | **Ignored for `generate` — output is capped at 1000 tokens regardless.** Use `chat` for longer outputs. |
| `temperature` | number | No | 0.7 | Creativity (0–1) |

**Fields that are accepted but NOT implemented:**
- `tools` — accepted, silently ignored; tool calling not supported
- `response_format` — accepted but silently ignored (verified 2026-08-10; JSON mode not enforced)
- `stream` — accepted but always returns full text; streaming not supported

**Vision IS supported** (corrected 2026-08-10): pass `messages[].content` as an OpenAI-style array
containing an `image_url` part (public URL or `data:` URI). Verified on `gpt-4o` — a red PNG returned
"Red". The earlier "not processed" note was wrong.

#### `chat` — Multi-turn conversation

```json
{
  "action": "chat",
  "messages": [
    {"role": "user", "content": "Write a caption"},
    {"role": "assistant", "content": "Here's a caption: ..."},
    {"role": "user", "content": "Make it shorter"}
  ],
  "model": "gpt-4o",
  "systemPrompt": "You are a copywriter.",
  "maxTokens": 4000
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | Must be `"chat"` |
| `messages` | array | Yes | Alternating `user`/`assistant` turns; last must be `role: "user"` |
| `messages[].role` | string | Yes | `"user"` or `"assistant"` |
| `messages[].content` | string or array | Yes | Text content |
| `model` | string | No | `gpt-4o-mini` — must be a valid OpenAI model ID |
| `systemPrompt` | string | No | Prepended to conversation |
| `maxTokens` | integer | No | 1200 — **honored by `chat`**; no hardcoded cap observed |
| `temperature` | number | No | 0.7 | Creativity (0–1) |

---

## Recipe: Long Structured Output

For any output that might exceed 1000 tokens (arc generation, show notes, research synthesis, multi-section structured content), use `chat` with a single user message wrapping the full prompt:

```javascript
const response = await fetch("https://audos.com/api/hooks/execute/workspace-351699/ai-api", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    action: "chat",
    model: "gpt-4.1",
    messages: [
      {
        role: "user",
        content: systemPrompt + "\n\n" + userPrompt   // pack system context here if needed
      }
    ],
    systemPrompt: "You are a podcast production assistant.",
    maxTokens: 8000
  })
});
const data = await response.json();
// data.usage.completionTokens will be close to your requested limit if output is long
// Always check: if completionTokens is exactly 1000, you may have used generate by mistake
```

If you see `completionTokens: 1000` consistently across different `maxTokens` values, you are hitting the `generate` cap. Switch to `chat`.

---

## Response Shape

### Success Response

```json
{
  "success": true,
  "text": "The generated output text",
  "model": "gpt-4o-mini-2024-07-18",
  "usage": {
    "promptTokens": 45,
    "completionTokens": 120,
    "totalTokens": 165
  },
  "_meta": {
    "success": true,
    "durationMs": 850,
    "logs": []
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | `true` even when model fails silently (for `generate`). Check `text` length. |
| `text` | string | Generated text. Empty string if model routing failed via `generate`. |
| `model` | string | Resolved model name (e.g., `gpt-4o-2024-07-18`). **Always log this field.** |
| `usage` | object | Token counts |
| `usage.promptTokens` | integer | Tokens in the prompt |
| `usage.completionTokens` | integer | Tokens in the response. If this is exactly 1000, you've hit the `generate` cap. |
| `usage.totalTokens` | integer | Sum |
| `_meta` | object | Audos metadata |
| `_meta.durationMs` | integer | Total request time in milliseconds |

### Error Response

```json
{
  "error": "Missing required field: prompt",
  "_meta": {
    "success": true,
    "durationMs": 26,
    "logs": []
  }
}
```

For `generate` action with a broken model name, no `error` field appears — you get `success: true` and empty `text`. Use `chat` to get the actual error.

| Error Case | Response | Notes |
|-----------|---------|-------|
| Missing required field | `error: "Missing required field: prompt"` | Returned before calling the provider |
| Broken model name via `generate` | `success: true`, `text: ""` | Silent failure — no error field |
| Broken model name via `chat` | `error: "OpenAI proxy error: 404 ..."` | Error surfaced correctly |
| Large system prompt | Accepts up to 50KB+ | No hard limit discovered |

---

## Capability Matrix per Model

### Vision / Multimodal

**Status:** NOT supported via the Audos hook. `image_url` and related fields are accepted but silently ignored.

### Tool Calling / Function Calling

**Status:** NOT supported. Hook accepts `tools` field but ignores it. Models respond in natural language.

### JSON Mode

| Model | JSON Mode | Notes |
|-------|-----------|-------|
| `gpt-4o`, `gpt-4o-mini`, `gpt-4.1`, `gpt-4-turbo` | Supported | Include `"response_format": {"type": "json_object"}` |
| All other models | Not applicable | They don't work via this hook |

### Maximum Token Support

| Action | Effective Cap | Notes |
|--------|---------------|-------|
| `generate` | **1000 tokens** | Hardcoded; `maxTokens` field is ignored |
| `chat` | Provider limit | Honored; verified up to 8000+ tokens returned |

### Streaming

**Status:** NOT SUPPORTED. Hook accepts `stream: true` but always returns full text.

---

## Latency Analysis

Note: latency measurements for non-OpenAI models reflect error response time, not generation time — these models are routing to OpenAI and failing. The numbers below are for working OpenAI models only.

| Model | Avg Latency | Range | Notes |
|-------|---|---|---|
| `gpt-4o` | ~763ms | 659–817ms | Moderate |
| `gpt-4o-mini` | ~1071ms | 952–1400ms | More variance |
| `gpt-4-turbo` | Slower than gpt-4o | — | Not re-tested post-proxy discovery |

---

## What's NOT Supported

| Feature | Status | Note |
|---------|--------|------|
| **Non-OpenAI models** | No | Proxy silently 404s; `generate` hides the error |
| **`generate` output > 1000 tokens** | No | Hardcoded cap; use `chat` instead |
| **Streaming** | No | Always returns full text |
| **Vision / Images** | No | Accepts `image_url` but doesn't process it |
| **Tool Calling** | No | Accepts `tools` field but ignores it |
| **Batch Mode** | Not applicable | Hook is single-request only |

---

## Authentication & Credentials

- **Public endpoint** — No API key required; workspace ID is embedded in the URL
- **Rate limiting** — No known limits; subject to provider-level rate limits
- **Cost** — Charges to Audos account at OpenAI rates; token counts available in response usage

---

## Error Handling

### Quick Reference

| Condition | Response | Fix |
|-----------|----------|-----|
| Missing `prompt` (for `generate`) | `error: "Missing required field: prompt"` | Add `"prompt": "..."` |
| Missing `messages` (for `chat`) | `error: "Missing required field: messages"` | Add `"messages": [...]` |
| Non-OpenAI model via `generate` | `success: true`, empty `text` | Switch to `chat` to see the actual error; use a gpt-* model |
| Non-OpenAI model via `chat` | `error: "OpenAI proxy error: 404 ..."` | Use a valid OpenAI model ID |
| Oversized payload | Hook accepts; provider may reject | Test with real request size |

### Example: Provider Error (visible via `chat` action)

```json
{
  "error": "OpenAI proxy error: 404 {\"error\": {\"message\": \"The model 'claude-sonnet-4-6' does not exist or you do not have access to it\"}}",
  "_meta": {
    "success": true,
    "durationMs": 321,
    "logs": [
      "[ERROR] OpenAI proxy error: 404 ..."
    ]
  }
}
```

---

## Known Issues & Workarounds

### Issue 1: Non-OpenAI model names accepted but broken

**Symptom:** Request succeeds, `text: ""`, no error (via `generate`). Or explicit 404 error (via `chat`).

**Root cause:** The hook is an OpenAI proxy. Non-OpenAI model IDs are passed to OpenAI's API which returns `model_not_found`. The `generate` action swallows this error.

**Workaround:** Use only `gpt-4.1`, `gpt-4o`, `gpt-4o-mini`, or `gpt-4-turbo`. Use `chat` action during development to catch model routing errors.

### Issue 2: `generate` output silently capped at 1000 tokens

**Symptom:** Structured output is consistently truncated. `completionTokens` is always 1000 regardless of the `maxTokens` value in the request.

**Root cause:** Hardcoded cap in the `generate` action handler. `maxTokens` parameter is ignored.

**Workaround:** Use `chat` action instead. Wrap the prompt as a single user message. `chat` honors `maxTokens`.

### Issue 3: Vision not exposed via hook

**Symptom:** `image_url` and related fields are accepted but silently ignored.

**Workaround:** Call provider APIs directly if vision is needed.

### Issue 4: Tool calling not exposed

**Symptom:** Hook accepts `tools` field; models respond in natural language.

**Workaround:** Call provider APIs directly; or prompt the model to describe tool calls as JSON.

### Issue 5: gpt-4o-mini is the default — consider being explicit

**Symptom:** If you omit `model`, you get gpt-4o-mini.

**Workaround:** Always specify `model` explicitly. Use `gpt-4.1` or `gpt-4o` for complex structured generation.

---

## Usage Examples

### Generate a short caption

```javascript
const response = await fetch("https://audos.com/api/hooks/execute/workspace-351699/ai-api", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    action: "generate",
    model: "gpt-4o",
    prompt: "Write a 2-sentence caption for a podcast about slow pandemics",
    systemPrompt: "Be poetic and thought-provoking",
    maxTokens: 100
  })
});
const data = await response.json();
console.log(data.text);
```

### Long structured output (use `chat`, not `generate`)

```javascript
const response = await fetch("https://audos.com/api/hooks/execute/workspace-351699/ai-api", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    action: "chat",
    model: "gpt-4.1",
    messages: [
      { role: "user", content: "Generate a full 4-section episode arc with transitions and a scorecard..." }
    ],
    systemPrompt: "You are a podcast production assistant.",
    maxTokens: 8000
  })
});
const data = await response.json();
// Check completionTokens — if it's 1000, something routed to generate instead
console.log(data.usage.completionTokens);
console.log(data.text);
```

### Iterative refinement (multi-turn)

```javascript
// Round 1
const r1 = await fetch("...", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    action: "chat",
    model: "gpt-4o",
    messages: [{ role: "user", content: "Write a caption" }],
    maxTokens: 100
  })
});
const round1 = await r1.json();

// Round 2 — user gives feedback
const r2 = await fetch("...", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    action: "chat",
    model: "gpt-4o",
    messages: [
      { role: "user", content: "Write a caption" },
      { role: "assistant", content: round1.text },
      { role: "user", content: "Make it one sentence only" }
    ],
    maxTokens: 100
  })
});
const round2 = await r2.json();
console.log(round2.text);
```

### JSON output

```javascript
const response = await fetch("...", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    action: "generate",
    model: "gpt-4o",
    prompt: "Return a JSON object with episode title, description, and 3 key themes",
    response_format: { type: "json_object" },
    maxTokens: 300
  })
});
const data = await response.json();
const parsed = JSON.parse(data.text);
```

---

## Billing & Cost

Token usage is reported in the `usage` field of every response. Charges are applied to the Audos workspace account at OpenAI rates (the hook is an OpenAI proxy).

```json
"usage": {
  "promptTokens": 45,
  "completionTokens": 120,
  "totalTokens": 165
}
```

---

## Last Probed

- **Date:** April 14, 2026
- **Models tested:** gpt-4o, gpt-4o-mini, gpt-4-turbo, gpt-4.1, claude-sonnet-4-6, claude-opus-4-6, claude-haiku-4-5-20251001, claude-sonnet-4-5, gpt-5, o1-preview, o3-mini, gemini-1.5-pro, gemini-2.0-flash, deepseek-chat, moonshotai/kimi-k2.5
- **Capabilities tested:** model routing, generate vs chat token caps (4 probes at 200/500/2000/8000 maxTokens), vision, tool calling, JSON mode, streaming, large system prompts, error surfacing by action type
- **Probe cost estimate:** ~$0.015

---

## See Also

- [AI Generation API](./throughline/ai-generation-api.md) — Quick reference for the two actions (`generate`, `chat`)
- [Throughline API Reference](./throughline/throughline-api-reference.md) — All endpoints
- Skill: [audos-platform](../skill/SKILL.md) — How to work with Audos APIs
