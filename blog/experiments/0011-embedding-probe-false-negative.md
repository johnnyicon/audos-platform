---
date: 2026-07-16
area: ai-api
status: corrected
label: "Embedding-generation probe, round 3: a false negative from guessing the wrong call shape"
---

**Hypothesis:** Does Audos expose a native embedding-generation primitive, the way it exposes
`platform.generateText`?

**Method:** Exhaustive live introspection of the `platform` global — every method and prototype level, a
regex sweep for embed/vector/encod/semantic — plus a sweep of `platform.integrations.isAvailable()` with
guessed provider-feature names like `openai-embeddings`, `vector-search`, `pinecone`.

**Result: confirmed absent — and wrong.** Every guessed name returned `false`; concluded no embedding
capability exists at all. This was corrected the same day by a second sub-job (see below) that tried the
actual gateway method instead of only guessing feature names — it never tried `platform.integrations.
proxy()` directly, and never tried the bare provider name `openai`. The lesson that outlasted the
specific finding: a capability search that only tries plausible-sounding names, without also trying the
platform's generic passthrough mechanism, will report false negatives for anything gated behind that
passthrough.

See `CHANGELOG.md`, "TEST E, corrected" section; `BACKLOG.md #14` (original text, corrected in place).
