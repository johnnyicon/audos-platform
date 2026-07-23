---
date: 2026-07-16
area: db-api
status: corrected
label: "Vector search validity, round 2: fake 5-float vectors don't prove semantic search works"
---

**Hypothesis:** With no pgvector enablement path available, is the fallback — JSON-array storage plus
brute-force cosine similarity — actually fast enough at realistic row counts to be a real answer, or a
scaling risk we haven't measured?

**Method:** Benchmarked brute-force cosine similarity over JSON-stored vectors at 50, 300, and 1,000
rows, inside a real hook against the live platform (job #82919) — using placeholder 5-float vectors, not
real embeddings.

**Result: valid but narrower than first claimed — a methodological gap we caught and corrected same
day.** Sub-millisecond throughout, 0.02ms to 0.28ms average. That genuinely answers "does row-count
scaling hold" — cosine arithmetic costs the same regardless of what the numbers mean. It does **not**
answer "does real semantic search work," since 5-float placeholder vectors aren't remotely the same
shape as a real ~1,536-float embedding. The conclusion was walked back the same day from "answered" to
"reopened, pending a real-dimensionality re-test" — see the round-3 re-benchmark below.

See `CHANGELOG.md`, "2026-07-16 — Capability test round 2" and the correction immediately below it;
`docs/platform/06-capabilities-reference.md`, Vector/Embedding Storage section.
