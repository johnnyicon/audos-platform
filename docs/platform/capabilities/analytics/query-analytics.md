# `query_analytics` — workspace overview metrics

**Status: ✅ verified live.** Otto-chat-triggered tool (see
`../../29-otto-tool-surface-vs-app-callable-hooks.md` — not a `platform.*` hook your app code can call).

**Params:** `metric` (visitors/conversions/sessions/overview), `days` (default 30), or explicit
`startDate`/`endDate`. Free.

**Verified 2026-07-23**, `metric=overview, days=30` against the real DoKnow workspace
(`8a65a4ac-5a22-435f-b55f-c41ea34ca00d`): **5 sessions, 4 contacts captured, 80.0% conversion rate.**

The contact count (4) was independently cross-checked against a *separate* tool call —
`query_data_source(sourceId=contacts)` (see `query-data-source.md`) — and matched exactly. That's a real
second-path corroboration, not the same number restated twice by the same tool.

Source: Otto chat, 2026-07-23. Narrative: `../../../../blog/experiments/0030-analytics-reporting-live-verification.md`.
