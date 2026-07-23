---
date: 2026-07-16
area: db-api
status: confirmed
label: "Vector dimensionality re-benchmark: does brute-force cosine hold at real embedding size?"
---

**Hypothesis:** Round 2's brute-force cosine similarity benchmark used 5-float placeholder vectors. Does
the same approach stay fast at realistic dimensionality — 1,536 floats, a real OpenAI embedding size —
and realistic row counts?

**Method:** Re-ran the identical benchmark from round 2 (50/300/1,000 rows), now with genuine
1,536-float vectors produced by the real embedding path confirmed above, three runs each for variance.

**Result: confirmed pass, closed.** Roughly 0.003ms/row scan time at 300–1,000 rows; a full 1,000-row
scan completed in 1–3ms. Roughly 307x more arithmetic per comparison barely moved the needle — wall-clock
is dominated by the database fetch itself (~250–540ms), not the comparison math. This closes the question
round 2 left open: brute-force cosine similarity is a real answer at DoKnow's scale, not just a
placeholder result.

See `docs/platform/06-capabilities-reference.md`, Vector/Embedding Storage results table; `BACKLOG.md
#13`.
