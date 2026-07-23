---
date: 2026-07-16
area: ai-api
status: confirmed
label: "Embedding-generation probe, round 3 corrected: proxy() actually works"
---

**Hypothesis:** Same question as the false negative above, tested through the platform's actual
integration gateway rather than by guessing feature names.

**Method:** Called `platform.integrations.isAvailable('openai')` with the bare provider name, then
`platform.integrations.proxy('openai', '/v1/embeddings', {method: 'POST', headers: {...}, body:
JSON.stringify({model: 'text-embedding-3-small', input: '...'})})`.

**Result: confirmed pass — corrects the false negative above.** `isAvailable('openai')` returned `true`;
the `proxy()` call succeeded and returned a real 1,536-float OpenAI embedding, with no API key supplied
by the calling code. `proxy` turned out to be a generic authenticated passthrough to a fixed allowlist
(`openai`, `stripe`, `twilio`, `heygen`) — Audos holds and injects the credential.

> A dangerous gotcha found in the process: `body` must be passed as a JSON string, not a plain object —
> passing an object silently returns the platform's own HTML index page with a 200 status, a false-success
> failure mode that looks identical to a real response until you check the content.

See `docs/platform/06-capabilities-reference.md`, Vector/Embedding Storage &rarr; "Native embedding
generation."
