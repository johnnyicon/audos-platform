---
date: 2026-07-14
area: app-build
status: corrected
label: "Ingestion pipeline probe, round 1: can DoKnow's pipeline be built on Audos as specced?"
---

**Hypothesis:** Per the platform's own (three-month-stale, never independently verified) capabilities
doc, DoKnow's full ingestion pipeline — upload, transcribe, chunk, embed, retrieval-grounded generation
— can be built entirely on Audos.

**Method:** A five-part disposable-hook probe: vector/embedding column and pgvector availability,
external `fetch()` reachability, file-upload size limits, scheduler/async hook firing, and
`platform.generateText`. The first attempt (job #82122) ran long with no cancel tool available and was
initially filed as possibly hung (`BACKLOG.md #10`) — later corrected: it wasn't hung, just genuinely
slow, completing at 45m01s after 67+ real tool actions. A concurrent retry (#82133) revealed a same-app
concurrency limit — it queued behind the first job rather than running in parallel.

**Result: mixed at the time, revised significantly by later rounds.** No native vector storage or
pgvector enablement path found; scheduled hooks confirmed to never fire (one schedule 2h07m overdue,
another 20m overdue, both `runCount: 0`). Upload, fetch, and `generateText` all worked cleanly. The
verdict at the time was "ingestion not buildable on Audos as originally specced" — later substantially
revised: embeddings turned out to be producible via an undocumented path (see the embedding-probe
experiments below), and the scheduler gap turned out to have a workable client-side substitute.

See `CHANGELOG.md`, "2026-07-15 — Capability test: ingestion pipeline verdict."
