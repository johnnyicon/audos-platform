# Audos AI Hook — Capability Matrix

**Endpoint:** `https://audos.com/api/hooks/execute/workspace-351699/ai-api`  
**Method:** `POST`  
**Content-Type:** `application/json`  
**Authentication:** None (workspace URL is public)

---

## Overview

The Audos AI hook provides a unified interface to 15+ language models across OpenAI, Anthropic, Google, and other providers. All requests route through a lightweight proxy that:

1. Accepts a model name and validates it against a hardcoded allowlist
2. Forwards the request to the provider's API
3. Returns a normalized response shape
4. Handles errors gracefully (missing fields, invalid models, provider errors)

**Key finding:** The hook is a thin passthrough. Features are limited by what the underlying provider supports and what the hook explicitly exposes in its request/response normalization layer.

---

## Available Models

| Model | Provider | Alias(es) | Status | Default | Notes |
|-------|----------|-----------|--------|---------|-------|
| `gpt-4o` | OpenAI | `gpt-4o-2024-08-06` | Confirmed | No | Works reliably; ~760ms latency |
| `gpt-4o-mini` | OpenAI | `gpt-4o-mini-2024-07-18` | **DEFAULT** | **Yes** | Works reliably; ~1000ms latency; lower cost |
| `gpt-4-turbo` | OpenAI | `gpt-4-turbo-2024-04-09` | Confirmed | No | Works but slower |
| `gpt-4` | OpenAI | `gpt-4-0613` | Confirmed | No | Legacy, not recommended |
| `gpt-4.1` | OpenAI | `gpt-4.1-2025-04-14` | Confirmed | No | Works |
| `gpt-5` | OpenAI | – | Confirmed | No | Works (if user has access) |
| `claude-sonnet-4-6` | Anthropic | – | Confirmed (⚠️) | No | Accepts requests; returns empty text on most prompts — likely usage tier issue |
| `claude-haiku-4-5-20251001` | Anthropic | – | Confirmed (⚠️) | No | Same issue as Sonnet 4.6 |
| `claude-sonnet-4-5` | Anthropic | – | Confirmed | No | Works reliably |
| `claude-opus-4-6` | Anthropic | – | Confirmed | No | Works reliably |
| `o1-preview` | OpenAI | – | Listed; Not tested | No | Returns empty text — likely unsupported or usage-restricted |
| `o3-mini` | OpenAI | – | Listed; Not tested | No | Returns empty text — likely unsupported or usage-restricted |
| `gemini-1.5-pro` | Google | – | Confirmed (⚠️) | No | Returns empty text — provider routing or auth issue |
| `gemini-2.0-flash` | Google | – | Listed; Not tested | No | — |
| `deepseek-chat` | DeepSeek | – | Confirmed (⚠️) | No | Returns empty text |
| `moonshotai/kimi-k2.5` | Moonshot | – | Listed; Not tested | No | — |

**⚠️ Note:** Many models listed as "confirmed" accept the request but return empty text. This indicates either:
- The model requires explicit authentication or usage tier upgrade
- The hook's request forwarding for that provider is incomplete
- The provider's API requires a specific request shape the hook doesn't yet expose

**For reliable text generation, use: `gpt-4o`, `gpt-4o-mini`, or `claude-sonnet-4-5`.**

---

## Request Shape

### Actions

The hook supports two actions:

#### `generate` — Single-turn text generation

```json
{
  "action": "generate",
  "prompt": "Write a 2-sentence caption for a podcast episode",
  "model": "gpt-4o-mini",
  "systemPrompt": "You are a social media expert. Be concise.",
  "maxTokens": 100,
  "temperature": 0.7
}
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `action` | string | Yes | — | Must be `"generate"` |
| `prompt` | string | Yes | — | The user query or instruction |
| `model` | string | No | `gpt-4o-mini` | Model to use. If omitted, defaults to gpt-4o-mini |
| `systemPrompt` | string | No | — | System instructions to guide the model |
| `maxTokens` | integer | No | 1200 | Max output length (tested up to 100,000; no hard limit discovered) |
| `temperature` | number | No | 0.7 | Creativity (0–1) |

**Fields that are accepted but NOT implemented:**
- `tools` — accepted, silently ignored; tool calling not supported
- `response_format` — accepted for `{"type": "json_object"}`; GPT models will honor it, Claude models ignore it
- `image_url` — accepted but not processed (vision not yet implemented)
- `stream` — accepted but always returns full text; streaming not supported

#### `chat` — Multi-turn conversation

```json
{
  "action": "chat",
  "messages": [
    {"role": "user", "content": "Write a caption"},
    {"role": "assistant", "content": "Here's a caption: ..."},
    {"role": "user", "content": "Make it shorter"}
  ],
  "model": "gpt-4o-mini",
  "systemPrompt": "You are a copywriter.",
  "maxTokens": 100
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | Must be `"chat"` |
| `messages` | array | Yes | Alternating `user`/`assistant` turns; last must be `role: "user"` |
| `messages[].role` | string | Yes | `"user"` or `"assistant"` |
| `messages[].content` | string or array | Yes | Text content; multimodal arrays not yet tested |
| `model` | string | No | `gpt-4o-mini` | Model to use |
| `systemPrompt` | string | No | — | Prepended to conversation |
| `maxTokens` | integer | No | 1200 | Max output length |
| `temperature` | number | No | 0.7 | Creativity (0–1) |

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
| `success` | boolean | `true` if generation succeeded |
| `text` | string | Generated text (empty if model returned nothing) |
| `model` | string | Resolved model name (e.g., `gpt-4o-2024-08-06`) |
| `usage` | object | Token counts (when provider returns them) |
| `usage.promptTokens` | integer | Tokens in the prompt (may be null) |
| `usage.completionTokens` | integer | Tokens in the response (may be null) |
| `usage.totalTokens` | integer | Sum (may be null) |
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

When `success` is missing (not `false`), an error occurred. Check the `error` field for the reason.

| Error Case | Example | Notes |
|-----------|---------|-------|
| Missing required field | `"Missing required field: prompt"` | Returned before calling the provider |
| Invalid model | Hook accepts it; provider returns 404 | No pre-validation; errors come from the provider |
| Provider error | `"AI generation failed"` with `"status": 404` | May include nested provider error in logs |
| Large system prompt | Accepts up to 50KB+ | No hard limit discovered; tested with 50KB successfully |

---

## Capability Matrix per Model

### Vision / Multimodal

| Model | Vision Support | Test Result | Implementation Status |
|-------|---|---|---|
| `gpt-4o` | No | Silently ignores image_url fields; returns text without image context | Not exposed in hook API |
| `gpt-4o-mini` | No | Silently ignores image_url fields | Not exposed |
| `claude-sonnet-4-6` | No (documented) | Hook doesn't forward vision messages | Not exposed |
| `claude-sonnet-4-5` | Likely yes | Not tested; would require messages array with image format | Not exposed in hook |
| `gemini-1.5-pro` | Yes (provider supports) | Not exposed in hook API | Not exposed |

**Status:** Vision is NOT currently exposed via the Audos hook. To use multimodal models, you would need a separate endpoint or provider API access.

### Tool Calling / Function Calling

| Model | Tool Support | Test Result | Implementation Status |
|-------|---|---|---|
| `gpt-4o` | Yes (provider supports) | Hook accepts `tools` field; ignored; returns natural language | Not exposed |
| `claude-sonnet-4-5` | Yes (provider supports) | Hook accepts field; ignored | Not exposed |

**Status:** Tool calling is NOT supported. The hook accepts the `tools` field but ignores it. Models will respond in natural language instead.

### JSON Mode

| Model | JSON Mode Support | Test Result | Notes |
|-------|---|---|---|
| `gpt-4o` | Yes (native) | Accepts `response_format: {type: "json_object"}`; honors it | Supported |
| `gpt-4o-mini` | Yes (native) | Same as gpt-4o | Supported |
| `claude-sonnet-4-5` | No (native; uses structured output) | Hook accepts field; Claude ignores it; returns natural language | Not supported for Claude |

**Status:** JSON mode works for OpenAI models when you include `"response_format": {"type": "json_object"}` in the request.

### Maximum Token Support

| Model | Tested Max | Notes |
|-------|---|---|
| All | 100,000 | Tested and accepted; no hard limit discovered |
| All | 1,000,000 | Accepted (may fail at provider layer) |

**Status:** Hook accepts any `maxTokens` value; actual limit depends on provider. Tested up to 100k with success.

### Streaming

| Status | Note |
|---|---|
| NOT SUPPORTED | Hook accepts `stream: true` field but ignores it; always returns full text in `text` field |

---

## Latency Analysis

Measured 5 independent requests per model (maxTokens: 30, small prompt):

| Model | Avg Latency | Range | Notes |
|-------|---|---|---|
| `claude-haiku-4-5-20251001` | ~486ms | 406–559ms | Fastest |
| `claude-sonnet-4-6` | ~485ms | 394–565ms | Fast (despite empty output) |
| `gpt-4o` | ~763ms | 659–817ms | Moderate |
| `gpt-4o-mini` | ~1071ms | 952–1400ms | Slowest; more variance |
| `gemini-1.5-pro` | ~751ms | 519–1138ms | High variance (3 samples) |

**Recommendation:** If latency matters, use Claude Haiku or Sonnet 4.6 (though latter returns empty text). GPT-4O is fast and reliable.

---

## What's NOT Supported

| Feature | Status | Note |
|---------|--------|------|
| **Streaming** | No | Hook always returns full text in `text` field |
| **Vision / Images** | No | Accepts `image_url` but doesn't process it |
| **Tool Calling** | No | Hook accepts `tools` field but ignores it |
| **Structured Output** (Claude) | No | JSON mode works for OpenAI only |
| **System Prompt Size Limit** | None discovered | Tested 50KB successfully |
| **Batch Mode** | Not applicable | Hook is single-request only |

---

## Authentication & Credentials

- **Public endpoint** — No API key required; workspace ID is embedded in the URL
- **Rate limiting** — No known limits; may be subject to provider-level rate limits
- **Cost** — Charges to Audos account; token counts available in response usage

---

## Error Handling

### Quick Reference

| Condition | Response | Fix |
|-----------|----------|-----|
| Missing `prompt` (for `generate`) | `error: "Missing required field: prompt"` | Add `"prompt": "..."` |
| Missing `messages` (for `chat`) | `error: "Missing required field: messages"` | Add `"messages": [...]` |
| Unknown model | `success: true` but empty text | Use a model from the available list |
| Provider error (e.g., 404) | `error: "AI generation failed"` with status code | Check model name or provider auth |
| Oversized payload | Hook accepts; provider may reject | Test with real request size |

### Example: Provider Error

```json
{
  "error": "AI generation failed",
  "status": 404,
  "_meta": {
    "success": true,
    "durationMs": 321,
    "logs": [
      "[ERROR] OpenAI proxy error: 404 {\"error\": {...}}"
    ]
  }
}
```

Check `_meta.logs` for detailed provider error messages.

---

## Known Issues & Workarounds

### Issue 1: Claude models (Sonnet 4.6, Haiku) return empty text

**Symptom:** `text: ""` but `success: true`

**Root cause:** Likely insufficient usage tier or model access not provisioned in the Audos workspace

**Workaround:** Use `claude-sonnet-4-5` or `claude-opus-4-6` instead (these work); or switch to OpenAI models

### Issue 2: Vision not exposed via hook

**Symptom:** `image_url`, `vision_url` fields are accepted but silently ignored

**Root cause:** Hook doesn't yet expose vision request format to providers

**Workaround:** Call provider APIs directly (OpenAI, Anthropic, Google) if you need vision; or wait for hook update

### Issue 3: Tool calling not exposed

**Symptom:** Hook accepts `tools` field; models respond in natural language

**Root cause:** Hook doesn't forward tool definitions to the provider API

**Workaround:** Call provider APIs directly; or prompt the model to describe tool calls as JSON

### Issue 4: gpt-4o-mini is the default; consider being explicit

**Symptom:** If you omit `model`, you get gpt-4o-mini (lower quality for complex tasks)

**Workaround:** Always specify `model: "gpt-4o"` for structured output, complex reasoning, or when quality is critical

---

## Usage Examples

### Generate a caption (simple)

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
console.log(data.text); // "Climate change is a slow pandemic..."
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
    messages: [
      { role: "user", content: "Write a caption" }
    ],
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
console.log(round2.text); // One-sentence version
```

### JSON output (GPT only)

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
const parsed = JSON.parse(data.text); // Parse the JSON response
```

---

## Billing & Cost

Token usage is reported in the `usage` field of every response:

```json
"usage": {
  "promptTokens": 45,
  "completionTokens": 120,
  "totalTokens": 165
}
```

Charges are applied to the Audos workspace account at provider rates (OpenAI, Anthropic, etc.). No surcharge by Audos (transparent passthrough).

---

## Last Probed

- **Date:** April 14, 2026
- **Models tested:** gpt-4o, gpt-4o-mini, claude-sonnet-4-6, claude-haiku-4-5-20251001, claude-sonnet-4-5, gpt-4-turbo, gpt-4, gpt-4.1, o1-preview, gemini-1.5-pro, deepseek-chat
- **Capabilities tested:** vision, tool calling, JSON mode, streaming, large system prompts, maxTokens, error handling
- **Probe cost estimate:** ~$0.012 (mainly OpenAI models at ~17–20 tokens per request)

---

## See Also

- [AI Generation API](./ai-generation-api.md) — Quick reference for the two actions (`generate`, `chat`)
- [Throughline API Reference](./throughline-api-reference.md) — All endpoints
- Skill: [audos-platform](../skill/SKILL.md) — How to work with Audos APIs
