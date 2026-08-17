---
date: 2026-08-14
product: DoKnow
status: pass
label: GPT-5 and the reasoning models were available on Audos the whole time, blocked by a single wrong parameter name in a hook nobody can read
---

# The frontier models were there all along

For months I believed the Audos AI hook could only reach GPT-4 class models. I had written it down as
settled: OpenAI proxy, `gpt-4o` and friends, no GPT-5, no reasoning models, no Claude. It was in the
capability matrix. I had repeated it in conversation.

It was wrong, and the way it was wrong is worth writing down, because the mistake was in how I read a
failure, not in the testing itself.

## The mistake

When you ask the hook for a model it can't serve, it fails. I had been reading "it failed" as "it isn't
available." Those are not the same thing, and the difference was sitting in the response the whole time.

The hook flattens every failure to the same unhelpful message:

```json
{ "error": "AI generation failed" }
```

That is what I had been looking at. But the response also carries a `_meta.logs` array, and the real
upstream reason survives in there untouched. Once I actually dumped the full payload instead of reading
the top-level error, two different failures separated out.

A model that genuinely isn't available returns this:

```
404 "The model `claude-sonnet-4-6` does not exist or you do not have access to it."
     code: model_not_found
```

But `gpt-5` returned something completely different:

```
400 "Unsupported parameter: 'max_tokens' is not supported with this model.
     Use 'max_completion_tokens' instead."
     code: unsupported_parameter
```

That is a parameter error, not a missing-model error. OpenAI only validates parameters *after* it has
resolved the model and authorised the key. So a 400 like that is positive proof the model exists and the
account can use it.

Every GPT-5 size, plus `o1`, `o3`, `o3-mini` and `o4-mini`, returned that same 400. Nine models, available
the entire time. The hook was simply asking for them with a field name OpenAI retired.

`o1-mini` returned a real 404 in the same sweep, which was the control that made the distinction
trustworthy rather than a hopeful reading.

## Why the hook couldn't just be fixed

The obvious move is to correct the hook. That turned out to be the harder problem.

The `ai-api` hook is a server function living in the workspace. I could not read its source. Neither could
Otto, which reported "Access denied" on every platform path it tried. It is not in the connected GitHub
repo either, and `AUDOS.md` lists what the repo syncs (components, React hooks, lib, data, apps, tools,
landing pages) with server functions absent from that list. Meanwhile the workspace's own server-hooks
tool says hooks "appear here when created through the agent," and the agent has no tool that writes them.

So: the agent says use the repo, the repo doesn't carry them, and the agent can't write them. For a builder
sitting in that workspace there is no path to the code at all.

## The way around

Editing the hook was impossible, but creating a new one was not. That capability had already been proven
back in July, when a Cursor delegation job created a temporary diagnostic hook to check whether pgvector
was installed. I had recorded that at the time and then failed to connect it to this problem, which cost
me more time than the bug did.

So instead of fixing `ai-api`, I built `ai2-api` alongside it. The whole fix is one conditional:

```js
if (/^(gpt-5|o[1-9])/.test(model)) {
  payload.max_completion_tokens = limit;
} else {
  payload.max_tokens = limit;
}
```

Everything else is the same call it always made, through
`platform.integrations.proxy('openai', '/v1/chat/completions', ...)`. Which matters, because that proxy
runs on **Audos's OpenAI key**. No key of my own, no BYOK. Frontier models in an app, hosted on Audos, on
their account.

Verified by calling it myself rather than trusting the job's report:

| Model | Served as | Result |
|---|---|---|
| `gpt-5` | `gpt-5-2025-08-07` | real output |
| `gpt-5-mini` | `gpt-5-mini-2025-08-07` | real output |
| `o3` | `o3-2025-04-16` | real output |
| `o4-mini` | `o4-mini-2025-04-16` | real output |
| `gpt-4o` (control) | `gpt-4o-2024-08-06` | real output |

## The second bug, which was quieter and worse

The first run through the fixed hook still failed on `gpt-5`, while `o3` worked fine. Same hook, same
parameter, different outcome.

Reasoning models spend part of the token budget thinking before they emit anything visible, and that spend
counts against the limit. At `maxTokens: 2000` the model used all 2000 tokens reasoning and returned an
empty string:

| maxTokens | completion tokens | reasoning tokens | visible text |
|---|---|---|---|
| 2000 | 2000 | 2000 | 0 characters |
| 4000 | 1746 | 1152 | 2266 characters |
| 8000 | 1926 | 1408 | 1892 characters |

This is nastier than the parameter bug. That one failed loudly. This one succeeds, bills you for 2000
tokens, and hands back nothing. A budget that is perfectly reasonable for GPT-4 silently produces empty
output on GPT-5.

The SDK client now defaults reasoning models to 4000 instead of 1200, and raises an explicit "raise
maxTokens" error rather than returning an empty string.

## What I take from this

**A generic error message is not a small problem.** Flattening every failure into "AI generation failed"
turned a one-line fix into months of believing a capability was missing. The information was always in the
response. It just wasn't where anyone would look.

**Distinguish "it failed" from "it isn't there."** I had evidence of the first and recorded the second.
The status code was doing real work and I was ignoring it in favour of the message.

**Check whether the wall has a door.** I had written off hook creation based on one tool being absent,
having already created a hook by another route weeks earlier. The blocker was in my notes, and so was the
answer.

The stock hook is still broken, and a builder still cannot read or edit their own workspace's server
function source. That part is filed. But the models were never missing.
