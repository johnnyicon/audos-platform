# `@audos/sdk` — AI client

Standalone, browser-safe access to a workspace's `ai-api` hook, so any app can use Audos's AI
(OpenAI text) over the API. Verified live 2026-08-10.

## Quick start

```ts
import { createAudosAI } from '@audos/sdk';

const ai = createAudosAI({ workspaceId: '351699' }); // no API key needed

const { text, usage } = await ai.complete('Write a one-line tagline for a knowledge app.');
console.log(text, usage); // usage: { promptTokens, completionTokens, totalTokens }
```

`complete()` is the method to reach for. It wraps your prompt as a chat turn, so it avoids the
`generate` action's 1000-token output cap.

### Multi-turn

```ts
const { text } = await ai.chat([
  { role: 'user', content: 'Draft a welcome email.' },
  { role: 'assistant', content: '…first draft…' },
  { role: 'user', content: 'Make it shorter and warmer.' },
], { model: 'gpt-4o', maxTokens: 800 });
```

`messages` must end with a `user` turn. `chat` honors `maxTokens` (unlike `generate`).

## Quick start (Go)

Same client, same behavior, in `sdk/go`:

```go
import audos "github.com/johnnyicon/audos-platform/sdk/go"

ai := audos.NewAI(audos.AIConfig{WorkspaceID: "351699"}) // no API key needed

r, err := ai.Complete("Write a one-line tagline for a goal-tracking app.", nil)
// r.Text, r.Usage.PromptTokens / CompletionTokens / TotalTokens

r2, _ := ai.Chat([]audos.ChatMessage{
    {Role: "user", Content: "Draft a nudge for a stalled goal."},
    {Role: "assistant", Content: "…first draft…"},
    {Role: "user", Content: "Warmer, one sentence."},
}, &audos.ChatOptions{Model: "gpt-4o", MaxTokens: 400})
```

Errors come back as `*audos.AIError` (`errors.As` to inspect `.Status` / `.ServerError`).

## Config

| Option | Default | Notes |
|---|---|---|
| `workspaceId` | — | Required. "351699" or "workspace-351699"; the prefix is stripped. |
| `baseUrl` | `https://audos.com` | |
| `hookName` | `ai-api` | The workspace hook to call. |
| `defaultModel` | `gpt-4o-mini` | |
| `hookSecret` | — | Sent as `x-hook-secret`. **Not** enforced by the reference hook — pass-through only, for a hook you build that checks it. Not authentication. |

## Vision (image understanding)

Works today on `gpt-4o` / `gpt-4.1` — verified 2026-08-10 (a red PNG came back "Red"). Pass an image
part, or use the `imageMessage` helper:

```ts
import { createAudosAI, imageMessage } from '@audos/sdk';
const ai = createAudosAI({ workspaceId: '351699' });
const { text } = await ai.chat(
  [imageMessage('What is in this image?', 'https://example.com/photo.jpg')], // URL or data: URI
  { model: 'gpt-4o' },
);
```

## What works, and what doesn't — from a 43-name live audit (2026-08-10)

The `AUDOS_AI_MODELS` export carries this list in code. Three tiers:

- **Working now:** the whole OpenAI **GPT-4/3.5 family** — `gpt-4o`, `gpt-4o-mini` (default), `gpt-4.1`,
  `gpt-4.1-mini`, `gpt-4.1-nano`, `gpt-4-turbo`, `gpt-4`, `gpt-3.5-turbo`. Vision on the 4o/4.1 families.
- **Frontier + reasoning models — WORKING as of 2026-08-14** via the corrected `ai2-api` hook:
  `gpt-5`, `gpt-5-mini`, `gpt-5-nano`, `o3`, `o4-mini` all verified returning real output. They need
  `max_completion_tokens` rather than `max_tokens`; the old `ai-api` hook sends the wrong one, so point
  the client at the corrected hook:

  ```ts
  const ai = createAudosAI({ workspaceId: '156396', hookName: 'ai2-api' });
  const { text } = await ai.complete('...', { model: 'gpt-5' });   // just works
  ```

  **Reasoning-token gotcha (handled for you).** These models spend part of the budget on *internal*
  reasoning before any visible text, and that counts against the limit. Measured: `gpt-5` on a moderately
  complex prompt burned all 2000 tokens reasoning and returned **empty text**; at 4000 it used ~1150
  reasoning plus real output. The client therefore defaults reasoning models to
  `REASONING_MIN_MAX_TOKENS` (4000) instead of 1200, and if the budget is still exhausted it throws a
  clear "raise maxTokens" error rather than a bare failure. If you set `maxTokens` yourself on a
  reasoning model, keep it generous.
- **Genuinely unavailable:** Claude, Fable, Gemini, `o1-mini`, `gpt-5.6-sol` — real `model_not_found`.
  Not OpenAI models, so no hook change reaches them through an OpenAI proxy. (Claude/Fable run on the
  Otto/job surface, which is not an app-callable API.)

- **Effort/reasoning params, `temperature`, `tools`, JSON-mode** — silently ignored by the reference hook.
- **Video / voice** are a different surface (Otto tools), not this hook.

## Two things to know before you ship it

1. **`ai-api` is a workspace-built server function, not a platform feature.** It exists in the
   Throughline workspace but 404s in others. If you point this client at a workspace that doesn't
   have the hook, `complete()`/`chat()` throw an `AudosAIError` saying exactly that. To use it in a
   new app you either call an existing workspace's hook, or build an `ai-api` hook in a workspace you
   control.

2. **The endpoint is public and CORS-open.** Convenient — you can call it straight from a browser,
   no proxy. But it means the workspace URL in your shipped client code is callable by anyone, and
   every call spends that workspace's OpenAI budget. For a public-facing app, prefer calling it from
   your own backend, or from a hook you control that adds its own rate limiting / secret check.

## Errors

Every non-success outcome throws `AudosAIError` with `.status` and `.serverError`. Notably,
`generate()` returns HTTP 200 with empty text when a model call fails — this client treats that empty
result as a thrown error rather than a silent blank. Prefer `complete()`/`chat()`.
