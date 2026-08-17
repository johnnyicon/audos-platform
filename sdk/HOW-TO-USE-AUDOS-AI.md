# How to use Audos AI

A complete, self-contained guide to calling Audos's AI from any application, in any language. Everything
here was verified live on 2026-08-14.

**What you get:** OpenAI text generation, including **GPT-5 and the o-series reasoning models**, plus
vision, running on **Audos's OpenAI account**. You do not need an OpenAI key of your own.

**What you don't get:** Claude, Gemini, or any non-OpenAI model. Also no streaming, tool calling, JSON
mode, or embeddings. See [Limits](#limits) for why.

---

## 1. The fastest possible start

It is one unauthenticated HTTPS POST. No SDK, no key, no setup:

```bash
curl -X POST https://audos.com/api/hooks/execute/workspace-{YOUR_WORKSPACE_ID}/ai2-api \
  -H "Content-Type: application/json" \
  -d '{
    "action": "chat",
    "model": "gpt-5",
    "messages": [{"role": "user", "content": "Explain vortex math in one sentence."}],
    "maxTokens": 4000
  }'
```

> **Substitute your own workspace ID.** The endpoint is unauthenticated, so whoever holds the URL spends
> that workspace's OpenAI budget. Note also that `ai2-api` is not a stock hook: it has to be created in
> your workspace first (see §9), so another workspace's ID won't work until you've done that.

Response:

```json
{
  "success": true,
  "text": "Vortex math is a numerological system that maps repeating digital-root patterns...",
  "model": "gpt-5-2025-08-07",
  "usage": { "prompt_tokens": 15, "completion_tokens": 185, "total_tokens": 200 }
}
```

That's the whole protocol. Everything below is detail.

---

## 2. The endpoint

```
POST https://audos.com/api/hooks/execute/workspace-{WORKSPACE_ID}/{HOOK_NAME}
Content-Type: application/json
```

| | |
|---|---|
| **Workspace** | `workspace-{YOUR_WORKSPACE_ID}`. Usage bills to that workspace. |
| **Hook** | **`ai2-api`** — use this one. |
| **Auth** | **None.** The workspace ID in the URL is the only identifier. |
| **CORS** | Open. The endpoint reflects any `Origin`, so browsers can call it directly with no proxy. |

### Use `ai2-api`, not `ai-api`

There are two hooks. `ai-api` is the original and **cannot reach GPT-5 or any reasoning model** — it sends
OpenAI's retired `max_tokens` parameter, which current models reject. `ai2-api` is the corrected version
and is what you want. If you get a 400 mentioning `max_completion_tokens`, you are on the wrong hook.

---

## 3. Request format

Two actions.

### `chat` — use this by default

```json
{
  "action": "chat",
  "model": "gpt-4o",
  "messages": [
    { "role": "user", "content": "Draft a welcome email." },
    { "role": "assistant", "content": "...first draft..." },
    { "role": "user", "content": "Shorter and warmer." }
  ],
  "systemPrompt": "You write concise product copy.",
  "maxTokens": 1200,
  "temperature": 0.7
}
```

| Field | Required | Notes |
|---|---|---|
| `action` | yes | `"chat"` |
| `messages` | yes | Alternating `user`/`assistant`. **Must end with a `user` turn.** |
| `model` | no | Defaults to `gpt-4o-mini`. See [Models](#4-models). |
| `systemPrompt` | no | Prepended as a system message. |
| `maxTokens` | no | Default 1200. **Reasoning models need ≥4000** — see [the trap](#6-the-reasoning-token-trap). |
| `temperature` | no | Ignored for reasoning models (they reject non-default values). |

### `generate` — single-turn, avoid it

```json
{ "action": "generate", "prompt": "Write a tagline.", "model": "gpt-4o" }
```

Works, but on the original hook it capped output at 1000 tokens regardless of what you asked for. **Prefer
`chat` with a single user message.** It has no cap and gives you the same result.

---

## 4. Models

### Working

| Model | Notes |
|---|---|
| `gpt-4o` | Fast, capable, vision |
| `gpt-4o-mini` | **Default.** Cheapest, vision |
| `gpt-4.1`, `gpt-4.1-mini`, `gpt-4.1-nano` | Vision |
| `gpt-4-turbo`, `gpt-4`, `gpt-3.5-turbo` | Older |
| **`gpt-5`, `gpt-5-mini`, `gpt-5-nano`** | Frontier. Needs `maxTokens` ≥4000 |
| **`o3`, `o3-mini`, `o4-mini`, `o1`** | Reasoning. Needs `maxTokens` ≥4000 |

### Not available

`claude-*` (any), `fable-5`, `gemini-*`, `deepseek-chat`, `o1-mini`, `chatgpt-4o-latest`. These return
`model_not_found` from OpenAI, because this is an OpenAI proxy and those aren't OpenAI models. There is no
setting that enables them.

> Claude *is* available on Audos, but only on the **build** surfaces (Audos Code, Cursor delegation jobs),
> not as an API your app can call. Getting Claude into an app would require your own Anthropic key via the
> platform's BYOK secrets mechanism, which is a different route entirely.

---

## 5. Vision

Pass `content` as an array instead of a string. Works on `gpt-4o` and `gpt-4.1`:

```json
{
  "action": "chat",
  "model": "gpt-4o",
  "maxTokens": 300,
  "messages": [{
    "role": "user",
    "content": [
      { "type": "text", "text": "What is in this image?" },
      { "type": "image_url", "image_url": { "url": "https://example.com/photo.jpg" } }
    ]
  }]
}
```

The URL can be public or a `data:image/png;base64,...` URI. Both verified working.

Two failure modes to expect, both reported clearly:
- **Unreachable URL** → `"Error while downloading <url>"`. OpenAI fetches the image itself, so the URL must
  be publicly reachable from their servers, not just from you.
- **Malformed image data** → `image_parse_error`. If you build a `data:` URI by hand, make sure it's a
  valid image; a truncated or hand-assembled PNG fails here rather than at the request layer.

---

## 6. The reasoning-token trap

**Read this before using GPT-5 or any o-series model.**

Reasoning models spend part of the token budget *thinking internally* before producing any visible text,
and that spend counts against your limit. If the budget runs out during reasoning, you get an **empty
string** and you are still billed.

Measured on `gpt-5` with a moderately complex prompt:

| `maxTokens` | completion tokens | of which reasoning | visible text |
|---|---|---|---|
| 2000 | 2000 | 2000 | **0 characters** |
| 4000 | 1746 | 1152 | 2266 characters |
| 8000 | 1926 | 1408 | 1892 characters |

**Rule: use `maxTokens` ≥ 4000 for `gpt-5*` and `o*` models, and higher for complex prompts.** A budget
that is perfectly fine for GPT-4 will silently return nothing on GPT-5.

You can detect it: the response has `finishReason: "length"` with empty `text`.

---

## 7. Errors

On failure the top-level `error` field is a generic string. **The real reason is in `_meta.logs`:**

```json
{
  "error": "AI generation failed",
  "status": 404,
  "_meta": {
    "logs": ["[ERROR] OpenAI proxy error: 404 {\"error\":{\"message\":\"The model `claude-sonnet-4-6` does not exist...\",\"code\":\"model_not_found\"}}"]
  }
}
```

Always read `_meta.logs[0]` when debugging. The `status` inside the body is the **upstream OpenAI** status:

| Upstream status | Meaning |
|---|---|
| `404` + `model_not_found` | That model genuinely isn't available |
| `400` + `unsupported_parameter` | You're on the old `ai-api` hook — switch to `ai2-api` |
| `finishReason: length`, empty text | Reasoning ate the budget — raise `maxTokens` |

**Also:** the `generate` action returns HTTP 200 with `text: ""` when the model call fails. Treat empty
text from `generate` as an error, never as an empty completion.

---

## 8. Using the SDK instead

If you're in TypeScript or Go, the SDK handles the hook choice, the reasoning-token default, error
extraction from `_meta.logs`, and the two different `usage` shapes.

**TypeScript** (browser-safe):

```ts
import { createAudosAI } from '@audos/sdk';

const ai = createAudosAI({ workspaceId: 'YOUR_WORKSPACE_ID', hookName: 'ai2-api' });
const { text, usage } = await ai.complete('Explain vortex math.', { model: 'gpt-5' });
```

**Go:**

```go
import audos "github.com/johnnyicon/audos-platform/sdk/go"

ai := audos.NewAI(audos.AIConfig{WorkspaceID: "YOUR_WORKSPACE_ID", HookName: "ai2-api"})
r, err := ai.Complete("Explain vortex math.", &audos.ChatOptions{Model: "gpt-5"})
```

Both default reasoning models to 4000 tokens automatically. See `AI-CLIENT.md` for the full API.

---

## Limits

**Not supported** (silently ignored if you send them): `tools` / function calling, `response_format` /
JSON mode, `stream`, and any reasoning-effort parameter. Embeddings are not reachable at all — the hook
only calls `/v1/chat/completions`.

**Security.** The endpoint is public and unauthenticated. Anyone with the URL can spend the workspace's
OpenAI budget. That's fine for a trusted internal system; for shipped client-side code, put your own
backend or a rate-limited hook in front of it.

**Billing.** All usage bills to the workspace in the URL, on Audos's OpenAI account.

**Stability.** `ai2-api` is a workspace-built server function, not a platform guarantee. It exists because
it was created deliberately; it is not present in other workspaces and Audos does not officially document
it. Treat it as durable-but-not-contractual.

---

## 9. Creating the `ai2-api` hook in your own workspace

`ai2-api` is not something Audos ships. It has to be created once per workspace, and the platform gives a
builder no direct way to write server functions: Otto has no hook-write tool, and the GitHub-connected repo
does not sync server functions. The route that does work is a **Cursor delegation job**, dispatched through
Otto, which can create hooks.

The hook source to give that job is in this repo at [`sdk/hooks/ai-api.reference.js`](hooks/ai-api.reference.js).
The essential part is branching the token parameter on the model id:

```js
if (/^(gpt-5|o[1-9])/.test(model)) {
  payload.max_completion_tokens = limit;   // GPT-5 + o-series
} else {
  payload.max_tokens = limit;              // GPT-4 and older
}
```

Everything else is a normal call through
`platform.integrations.proxy('openai', '/v1/chat/completions', { method:'POST', headers:{...},
body: JSON.stringify(payload) })`.

Two gotchas when authoring it: the proxy's `body` **must** be a JSON string (passing an object silently
returns the platform's HTML index page with status 200), and you should not forward `temperature` to
reasoning models, which reject non-default values.
